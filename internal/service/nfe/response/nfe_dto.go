package response

import (
	"encoding/xml"
)

// Estruturas mínimas para parsing

type RetEnviNFe struct {
	XMLName xml.Name `xml:"retEnviNFe"`
	ProtNFe ProtNFe  `xml:"protNFe"`
}

type ProtNFe struct {
	XMLName xml.Name `xml:"protNFe"`
	Inner   string   `xml:",innerxml"`
}

// Wrapper SOAP (para pegar o conteúdo dentro)
type SoapEnvelope struct {
	Body SoapBody `xml:"Body"`
}

type SoapBody struct {
	NfeResultMsg NfeResultMsg `xml:"nfeResultMsg"`
}

type NfeResultMsg struct {
	RetEnviNFe RetEnviNFe `xml:"retEnviNFe"`
}
