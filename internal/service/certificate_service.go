package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"github.com/beevik/etree"
	pkcs12 "software.sslmate.com/src/go-pkcs12"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func LoadCertFromS3(bucket, key, password string) (tls.Certificate, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return tls.Certificate{}, err
	}

	client := s3.NewFromConfig(cfg)

	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return tls.Certificate{}, err
	}
	defer resp.Body.Close()

	pfxData, err := io.ReadAll(resp.Body)
	if err != nil {
		return tls.Certificate{}, err
	}

	// 🔥 EXTRAÇÃO CORRETA
	privateKey, certificate, caCerts, err := pkcs12.DecodeChain(pfxData, password)
	if err != nil {
		return tls.Certificate{}, err
	}

	// 🔥 monta cadeia correta
	cert := tls.Certificate{
		PrivateKey:  privateKey,
		Certificate: make([][]byte, 0),
	}

	// 🔥 PRIMEIRO = leaf (OBRIGATÓRIO)
	cert.Certificate = append(cert.Certificate, certificate.Raw)

	// 🔥 depois cadeia (opcional, mas ok)
	for _, ca := range caCerts {
		cert.Certificate = append(cert.Certificate, ca.Raw)
	}

	// 🔥 parse do leaf
	cert.Leaf = certificate

	return cert, nil
}

func BuildEnviNFe(signedXML []byte, loteID string, indSinc string) (*etree.Element, error) {

	envi := etree.NewElement("enviNFe")
	envi.CreateAttr("xmlns", "http://www.portalfiscal.inf.br/nfe")
	envi.CreateAttr("versao", "4.00")

	// =========================
	// Lote
	// =========================
	envi.CreateElement("idLote").SetText(loteID)
	envi.CreateElement("indSinc").SetText(indSinc)

	// =========================
	// NF-e assinada (XML já pronto do xmlsec1)
	// =========================
	if err := injectRawXML(envi, signedXML); err != nil {
		return nil, err
	}

	return envi, nil
}

func injectRawXML(parent *etree.Element, rawXML []byte) error {

	doc := etree.NewDocument()

	if err := doc.ReadFromBytes(rawXML); err != nil {
		return err
	}

	root := doc.Root()
	if root == nil {
		return fmt.Errorf("XML inválido: sem root")
	}

	// 🔥 IMPORTANTE: adicionar como elemento, não texto
	parent.AddChild(root)

	return nil
}

func BuildSOAPEnvelope(enviNFe *etree.Element) ([]byte, error) {

	envelope := etree.NewDocument()

	// XML header
	envelope.CreateProcInst("xml", `version="1.0" encoding="utf-8"`)

	// =========================
	// SOAP Envelope
	// =========================
	env := envelope.CreateElement("soap12:Envelope")
	env.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	env.CreateAttr("xmlns:xsd", "http://www.w3.org/2001/XMLSchema")
	env.CreateAttr("xmlns:soap12", "http://www.w3.org/2003/05/soap-envelope")

	// =========================
	// Header
	// =========================
	header := env.CreateElement("soap12:Header")
	cabec := header.CreateElement("nfeCabecMsg")
	cabec.CreateAttr("xmlns", "http://www.portalfiscal.inf.br/nfe/wsdl/NFeAutorizacao4")

	cabec.CreateElement("versaoDados").SetText("4.00")
	cabec.CreateElement("cUF").SetText("35")

	// =========================
	// Body
	// =========================
	body := env.CreateElement("soap12:Body")
	dados := body.CreateElement("nfeDadosMsg")
	dados.CreateAttr("xmlns", "http://www.portalfiscal.inf.br/nfe/wsdl/NFeAutorizacao4")

	// =========================
	// 🔥 INJEÇÃO DO ENVI NFe (já estruturado)
	// =========================
	dados.AddChild(enviNFe)

	// =========================
	// Serialização final
	// =========================
	out, err := envelope.WriteToBytes()
	if err != nil {
		return nil, err
	}

	return out, nil
}
