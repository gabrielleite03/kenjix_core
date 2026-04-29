package nfe

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	NFeResultMsg NFeResultMsg `xml:"nfeResultMsg"`
}

type NFeResultMsg struct {
	RetEnviNFe RetEnviNFe `xml:"retEnviNFe"`
}

type RetEnviNFe struct {
	TpAmb    string   `xml:"tpAmb"`
	VerAplic string   `xml:"verAplic"`
	CStat    int      `xml:"cStat"`
	XMotivo  string   `xml:"xMotivo"`
	CUF      string   `xml:"cUF"`
	DhRecbto string   `xml:"dhRecbto"`
	ProtNFe  *ProtNFe `xml:"protNFe"`
}

type ProtNFe struct {
	InfProt InfProt `xml:"infProt"`
}

type InfProt struct {
	TpAmb    string `xml:"tpAmb"`
	VerAplic string `xml:"verAplic"`
	ChNFe    string `xml:"chNFe"`
	DhRecbto string `xml:"dhRecbto"`
	CStat    int    `xml:"cStat"`
	XMotivo  string `xml:"xMotivo"`
}
