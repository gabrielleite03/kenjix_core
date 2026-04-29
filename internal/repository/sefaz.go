// sefaz/client.go
package repository

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

func SignNFeXML(xmlData []byte, cert tls.Certificate) ([]byte, error) {

	// garante leaf
	if cert.Leaf == nil {
		var err error
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return nil, err
		}
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return nil, err
	}

	nfe := doc.FindElement(".//NFe")
	if nfe == nil {
		return nil, fmt.Errorf("NFe não encontrada")
	}

	infNFe := nfe.FindElement("infNFe")
	if infNFe == nil {
		return nil, fmt.Errorf("infNFe não encontrada")
	}

	id := infNFe.SelectAttrValue("Id", "")
	if id == "" {
		return nil, fmt.Errorf("Id não encontrada")
	}

	// garante Id correto
	infNFe.RemoveAttr("Id")
	infNFe.CreateAttr("Id", id)

	// contexto assinatura
	ctx := dsig.NewDefaultSigningContext(dsig.TLSCertKeyStore(cert))
	ctx.Hash = crypto.SHA1
	ctx.Canonicalizer = dsig.MakeC14N10RecCanonicalizer()
	ctx.IdAttribute = "Id"

	// 🔥 assinatura
	signed, err := ctx.SignEnveloped(infNFe)
	if err != nil {
		return nil, err
	}

	// pega signature
	signature := signed.FindElement("./Signature")
	if signature == nil {
		return nil, fmt.Errorf("Signature não encontrada")
	}

	// move assinatura para fora do infNFe (padrão NF-e)
	signed.RemoveChild(signature)
	nfe.AddChild(signature)

	// serializa
	return doc.WriteToBytes()
}

func SendNFe(xml []byte, cert tls.Certificate, endpoint string) ([]byte, error) {

	// 1. Assinar XML
	signedXML, err := SignNFeXML(xml, cert)
	if err != nil {
		return nil, err
	}

	// 2. RootCAs seguro
	roots, err := x509.SystemCertPool()
	if err != nil || roots == nil {
		return nil, fmt.Errorf("cannot load system cert pool")
	}

	// 3. TLS config limpo
	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,
		ServerName:   "homologacao.nfe.fazenda.sp.gov.br",
	}

	transport := &http.Transport{
		TLSClientConfig:   tlsConfig,
		ForceAttemptHTTP2: false, // SEFAZ às vezes é sensível
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	// 4. SOAP request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(signedXML))
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		`application/soap+xml; charset=utf-8; action="http://www.portalfiscal.inf.br/nfe/wsdl/NFeAutorizacao4/nfeAutorizacaoLote"`,
	)

	// 5. envio direto
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func SendNFeOLd(xml []byte, cert tls.Certificate, endpoint string) ([]byte, error) {

	// 🔥 SOAP 1.2 (CORRETO)
	envelope := xml

	// ===============================
	// 🔐 ROOT CAs
	// ===============================
	roots, err := x509.SystemCertPool()
	if err != nil || roots == nil {
		log.Println("⚠️ SystemCertPool vazio, criando novo pool")
		roots = x509.NewCertPool()
	}

	// linux
	/* koto remover para usar no linux
	roots = x509.NewCertPool()

	certs, _ := os.ReadFile("/etc/ssl/certs/ca-certificates.crt")
	roots.AppendCertsFromPEM(certs)
	*/
	// ===============================
	// 🔐 DEBUG CERTIFICADO
	// ===============================
	log.Println("🔐 Cert chain length:", len(cert.Certificate))

	// ===============================
	// 🔐 TLS CONFIG (COM RENEGOTIATION)
	// ===============================
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,

		Certificates: []tls.Certificate{cert},

		RootCAs: roots,

		// 🔥 ESSENCIAL
		ClientAuth: tls.RequireAndVerifyClientCert,

		// 🔥 evita problemas com cadeia
		InsecureSkipVerify: false,

		ServerName: "homologacao.nfe.fazenda.sp.gov.br",

		Renegotiation: tls.RenegotiateFreelyAsClient,
	}

	// ===============================
	// 🌐 HTTP CLIENT
	// ===============================
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// ===============================
	// 🧪 TESTE TLS
	// ===============================
	log.Println("🌐 Testando handshake TLS...")

	resp, err := client.Get("https://homologacao.nfe.fazenda.sp.gov.br")
	if err != nil {
		logTLSError("GET", err)
		return nil, err
	}
	resp.Body.Close()

	//log.Println("✅ TLS OK - status:", resp.Status)

	// ===============================
	// 🚀 POST SOAP
	// ===============================
	log.Println("📤 Enviando NF-e para SEFAZ...")
	log.Println(string(envelope))

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(envelope))
	if err != nil {
		return nil, err
	}

	// ✅ SOAP 1.2 correto
	req.Header.Set(
		"Content-Type",
		`application/soap+xml; charset=utf-8; action="http://www.portalfiscal.inf.br/nfe/wsdl/NFeAutorizacao4/nfeAutorizacaoLote"`)

	// 🔥 IMPORTANTE: alguns servidores quebram se isso existir duplicado
	req.Header.Del("SOAPAction")

	// opcional
	req.Header.Set("Accept", "application/soap+xml")

	// opcional (alguns ambientes exigem)
	req.Header.Set("SOAPAction", "http://www.portalfiscal.inf.br/nfe/wsdl/NFeAutorizacao4/nfeAutorizacaoLote")

	resp, err = client.Do(req)
	if err != nil {
		logTLSError("POST", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("📥 Resposta SEFAZ status:", resp.Status)
	//log.Println("📥 BODY:")
	log.Println(string(body))

	return body, nil
}

// ===============================
// 🔍 DEBUG TLS
// ===============================
func logTLSError(stage string, err error) {
	log.Printf("❌ %s ERROR: %+v\n", stage, err)

	if urlErr, ok := err.(*url.Error); ok {
		log.Printf("🔎 %s URL ERROR: %+v\n", stage, urlErr)

		if opErr, ok := urlErr.Err.(*net.OpError); ok {
			log.Printf("🔎 %s NET ERROR: %+v\n", stage, opErr)

			if opErr.Err != nil {
				log.Printf("🔎 %s INNER ERROR: %+v\n", stage, opErr.Err)
			}
		}
	}
}
