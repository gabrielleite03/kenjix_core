// service/nfe_service.go
package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabrielleite03/kenjix_core/internal/repository"
)

func EmitirNFeKoto(xmlData []byte) ([]byte, error) {
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
	randomCNF := fmt.Sprintf("%08d", rand.Intn(99999999))
	fileName := fmt.Sprintf("%s-nfe.xml", randomCNF)

	certificateFileName := "clean.pfx"

	rawPath := filepath.Join(rawDir, fileName)
	signedPath := filepath.Join(signedDir, fileName)
	soapPath := filepath.Join(soapDir, fileName)
	certificatePath := filepath.Join(certificateDir, certificateFileName)

	err := os.WriteFile(rawPath, xmlData, 0644)
	if err != nil {
		fmt.Println("Erro:", err.Error())
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
		fmt.Println("Erro:", err.Error())
		return nil, err
	}

	cmd = exec.Command(
		"xmlsec1",
		"--verify",
		"--pkcs12", certificatePath,
		"--pwd", certificatePassword,
		"--id-attr:Id", "infNFe",
		"signed.xml",
	)

	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Erro:", err)
		fmt.Println(string(out))
		return nil, err
	}

	fmt.Println("OK assinatura verificada")

	signedXML, _ := os.ReadFile(signedPath)

	enviXML, _ := BuildEnviNFe(signedXML, "1", "1")
	soapEnvelop, err := BuildSOAPEnvelopeKoto(enviXML)
	os.WriteFile(soapPath, soapEnvelop, 0644)
	if err != nil {
		return nil, err
	}

	cert, err := LoadCertFromS3("aws-s3-site-kejipet", "certs/kenjipet.pfx", "K0t0net#g4")
	resp, err := repository.SendNFeOLd(
		soapEnvelop,
		/* ainda precisa do cert para TLS */
		/* pode manter seu LoadCertFromS3 aqui */
		cert,
		"https://homologacao.nfe.fazenda.sp.gov.br/ws/nfeautorizacao4.asmx",
	)

	return resp, err

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
