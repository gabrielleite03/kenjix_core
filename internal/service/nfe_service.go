// service/nfe_service.go
package service

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabrielleite03/kenjix_core/internal/repository"
	"github.com/gabrielleite03/kenjix_core/internal/service/nfe"
	"github.com/gabrielleite03/kenjix_core/internal/service/nfe/response"
)

func EmitirNFe(xmlData []byte) ([]byte, error) {
	certificatePassword := os.Getenv("CERT_PASSWORD")
	baseDir := filepath.Join(os.TempDir(), "nfe")

	rawDir := filepath.Join(baseDir, "raw")
	signedDir := filepath.Join(baseDir, "signed")
	soapDir := filepath.Join(baseDir, "soap")
	certificateDir := filepath.Join(baseDir, "certificate")

	os.MkdirAll(rawDir, 0755)
	os.MkdirAll(signedDir, 0755)
	os.MkdirAll(soapDir, 0755)
	os.MkdirAll(certificateDir, 0755)
	randomCNF := fmt.Sprintf("%10d", rand.Intn(99999999))
	fileName := fmt.Sprintf("%s-nfe.xml", randomCNF)

	certificateFileName := "clean.pfx"

	rawPath := filepath.Join(rawDir, fileName)
	signedPath := filepath.Join(signedDir, fileName)
	soapPath := filepath.Join(soapDir, fileName)
	certificatePath := filepath.Join(certificateDir, certificateFileName)

	err := os.WriteFile(rawPath, xmlData, 0644)
	if err != nil {
		fmt.Println("Error to write file:", err.Error())
		return nil, err
	}

	certificatePath, err = saveCertPFXFromS3(
		"aws-s3-site-kejipet",
		"certs/"+certificateFileName,
		certificatePath,
	)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(
		"xmlsec1",
		"--sign",
		"--output", signedPath,
		"--pkcs12", certificatePath,
		"--pwd", certificatePassword,
		"--id-attr:Id", "infNFe",
		"--node-xpath", "//*[local-name()='infNFe']",
		rawPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error to sign file:", err.Error())
		return nil, err
	}

	cmd = exec.Command(
		"xmlsec1",
		"--verify",
		"--pkcs12", certificatePath,
		"--pwd", certificatePassword,
		"--id-attr:Id", "infNFe",
		signedPath,
	)

	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error to verify signature:", err.Error())
		fmt.Println(string(out))
		return nil, err
	}

	fmt.Println("OK assinatura verificada")

	signedXML, _ := os.ReadFile(signedPath)

	enviXML, _ := BuildEnviNFe(signedXML, "1", "1")
	soapEnvelop, err := BuildSOAPEnvelope(enviXML)
	os.WriteFile(soapPath, soapEnvelop, 0644)
	if err != nil {
		return nil, err
	}

	cert, err := LoadCertFromS3("aws-s3-site-kejipet", "certs/kenjipet.pfx", certificatePassword)
	resp, err := repository.SendNFeOLd(
		soapEnvelop,
		/* ainda precisa do cert para TLS */
		/* pode manter seu LoadCertFromS3 aqui */
		cert,
		"https://homologacao.nfe.fazenda.sp.gov.br/ws/nfeautorizacao4.asmx",
	)

	ret, err := parseNFeResponse(resp)
	if err != nil {
		return nil, err
	}

	log.Println("📊 Status lote:", ret.CStat, ret.XMotivo)

	if ret.ProtNFe != nil {
		log.Println("📊 Status NF:", ret.ProtNFe.InfProt.CStat)
		log.Println("📊 Motivo:", ret.ProtNFe.InfProt.XMotivo)
		log.Println("📊 Chave:", ret.ProtNFe.InfProt.ChNFe)
		uploadXMLToS3(soapEnvelop, fmt.Sprintf("%s.xml", ret.ProtNFe.InfProt.ChNFe), "uploads/nfe/sended/")
		uploadXMLToS3(resp, fmt.Sprintf("%s.xml", ret.ProtNFe.InfProt.ChNFe), "uploads/nfe/return/")

		nfeProcXML, err := BuildNFeProc(soapEnvelop, resp)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(nfeProcXML))
	}

	/*

			cStat	Significado	Precisa consultar recibo?
		103	Lote recebido	✅ SIM
		105	Em processamento	✅ SIM
		104	Lote processado	❌ NÃO
			if resp.cStat == 103 {
		    nRec := resp.nRec

		    // 2. consulta até finalizar
		    for {
		        ret := ConsultarRecibo(nRec)

		        if ret.cStat == 105 {
		            time.Sleep(2 * time.Second)
		            continue
		        }

		        if ret.cStat == 100 {
		            // sucesso
		            salvarXML(ret)
		            gerarDANFE(ret)
		            break
		        }

		        // erro
		        log.Println("Rejeição:", ret.xMotivo)
		        break
		    }
		}

	*/

	return resp, err

}

func uploadXMLToS3(xml []byte, fileName string, urlS3 string) error {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("aws-s3-site-kejipet"),
		Key:         aws.String(urlS3 + fileName),
		Body:        bytes.NewReader(xml),
		ContentType: aws.String("application/xml"), // 👈 importante
	})

	return err
}

func BuildNFeProc(nfeXML []byte, retXML []byte) ([]byte, error) {
	// =========================
	// 1. Extrair protNFe do retorno
	// =========================
	var soap response.SoapEnvelope
	if err := xml.Unmarshal(retXML, &soap); err != nil {
		return nil, err
	}

	prot := soap.Body.NfeResultMsg.RetEnviNFe.ProtNFe
	if prot.Inner == "" {
		return nil, errors.New("protNFe não encontrado no retorno")
	}

	// =========================
	// 2. Extrair NFe (sem envelope)
	// =========================
	nfeContent, err := extractNFe(nfeXML)
	if err != nil {
		return nil, err
	}

	// =========================
	// 3. Montar nfeProc
	// =========================
	var buffer bytes.Buffer

	buffer.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buffer.WriteString(`<nfeProc xmlns="http://www.portalfiscal.inf.br/nfe" versao="4.00">`)

	buffer.Write(nfeContent)

	buffer.WriteString(`<protNFe versao="4.00">`)
	buffer.WriteString(prot.Inner)
	buffer.WriteString(`</protNFe>`)

	buffer.WriteString(`</nfeProc>`)

	return buffer.Bytes(), nil
}

func extractNFe(xmlData []byte) ([]byte, error) {
	start := bytes.Index(xmlData, []byte("<NFe"))
	end := bytes.Index(xmlData, []byte("</NFe>"))

	if start == -1 || end == -1 {
		return nil, errors.New("tag <NFe> não encontrada")
	}

	end += len("</NFe>")

	return xmlData[start:end], nil
}

func parseNFeResponse(data []byte) (*nfe.RetEnviNFe, error) {
	var envelope nfe.Envelope

	err := xml.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}

	return &envelope.Body.NFeResultMsg.RetEnviNFe, nil
}

func saveCertPFXFromS3(bucket, key, outputPath string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	pfxData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// garante diretório
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", err
	}

	// salva arquivo .pfx
	if err := os.WriteFile(outputPath, pfxData, 0644); err != nil {
		return "", err
	}

	return outputPath, nil
}

func GenerateDanfe(xmlPath, pdfPath string) error {
	cmd := exec.Command("php", "gerar_danfe.php", xmlPath, pdfPath)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro DANFE: %s - %s", err, string(out))
	}

	return nil
}
