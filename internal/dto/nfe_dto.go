// model/nfe.go
package dto

import "encoding/xml"

type EnviNFe struct {
	XMLName xml.Name `xml:"enviNFe"`
	Xmlns   string   `xml:"xmlns,attr"`
	Versao  string   `xml:"versao,attr"`
	IdLote  string   `xml:"idLote"`
	IndSinc string   `xml:"indSinc"`
	TpAmb   string   `xml:"tpAmb"`

	NFe NFe `xml:"NFe"`
}

type NFe struct {
	XMLName xml.Name `xml:"NFe"`
	Xmlns   string   `xml:"xmlns,attr"`
	InfNFe  InfNFe   `xml:"infNFe"`
}

type InfNFe struct {
	ID     string `xml:"Id,attr"`
	Versao string `xml:"versao,attr"`

	Ide   Ide   `xml:"ide"`
	Emit  Emit  `xml:"emit"`
	Dest  Dest  `xml:"dest"`
	Det   []Det `xml:"det"`
	Total Total `xml:"total"`

	Transp Transp `xml:"transp"`
	Pag    Pag    `xml:"pag"`
}

type Ide struct {
	CUF    string  `xml:"cUF"`
	CNF    string  `xml:"cNF"`
	NatOp  string  `xml:"natOp"`
	Mod    string  `xml:"mod"`
	Serie  string  `xml:"serie"`
	NNF    string  `xml:"nNF"`
	DhEmi  string  `xml:"dhEmi"`
	TpNF   string  `xml:"tpNF"`
	IdDest string  `xml:"idDest"`
	CMunFG string  `xml:"cMunFG"`
	TpImp  string  `xml:"tpImp"`
	TpEmis string  `xml:"tpEmis"`
	CDV    *string `xml:"cDV,omitempty"`

	TpAmb    string `xml:"tpAmb"` // 👈 TEM QUE ESTAR AQUI
	FinNFe   string `xml:"finNFe"`
	IndFinal string `xml:"indFinal"`
	IndPres  string `xml:"indPres"`
	ProcEmi  string `xml:"procEmi"`
	VerProc  string `xml:"verProc"`
}

type Emit struct {
	CNPJ      string `xml:"CNPJ,omitempty"`
	CPF       string `xml:"CPF,omitempty"`
	XNome     string `xml:"xNome"`
	XFant     string `xml:"xFant,omitempty"`
	EnderEmit Ender  `xml:"enderEmit"`

	IE   string `xml:"IE,omitempty"`
	IEST string `xml:"IEST,omitempty"`
	IM   string `xml:"IM,omitempty"`
	CNAE string `xml:"CNAE,omitempty"`

	CRT string `xml:"CRT"`
}

type Dest struct {
	CNPJ string `xml:"CNPJ,omitempty"`
	CPF  string `xml:"CPF,omitempty"`
	IE   string `xml:"IE,omitempty"`
	ISUF string `xml:"ISUF,omitempty"`
	IM   string `xml:"IM,omitempty"`

	XNome     string `xml:"xNome"`
	EnderDest Ender  `xml:"enderDest"`

	IndIEDest string `xml:"indIEDest,omitempty"`
	Email     string `xml:"email,omitempty"`
}

type Ender struct {
	XLgr    string `xml:"xLgr"`
	Nro     string `xml:"nro"`
	XCpl    string `xml:"xCpl,omitempty"`
	XBairro string `xml:"xBairro"`
	CMun    string `xml:"cMun"`
	XMun    string `xml:"xMun"`
	UF      string `xml:"UF"`
	CEP     string `xml:"CEP"`
	CPais   string `xml:"cPais"`
	XPais   string `xml:"xPais"`
	Fone    string `xml:"fone,omitempty"`
}

type Det struct {
	NItem   int     `xml:"nItem,attr"`
	Prod    Prod    `xml:"prod"`
	Imposto Imposto `xml:"imposto"`
}

type Prod struct {
	CProd  string `xml:"cProd"`
	CEAN   string `xml:"cEAN,omitempty"`
	XProd  string `xml:"xProd"`
	NCM    string `xml:"NCM"`
	CFOP   string `xml:"CFOP"`
	UCom   string `xml:"uCom"`
	QCom   string `xml:"qCom"`
	VUnCom string `xml:"vUnCom"`
	VProd  string `xml:"vProd"`
}

type Imposto struct {
	ICMS   ICMS   `xml:"ICMS"`
	PIS    PIS    `xml:"PIS"`
	COFINS COFINS `xml:"COFINS"`
}

type ICMS struct {
	ICMSSN102 ICMSSN102 `xml:"ICMSSN102"`
}

type ICMSSN102 struct {
	Orig  string `xml:"orig"`
	CSOSN string `xml:"CSOSN"`
}

type Total struct {
	ICMSTot ICMSTot `xml:"ICMSTot"`
}

type ICMSTot struct {
	VBC        string `xml:"vBC"`
	VICMS      string `xml:"vICMS"`
	VICMSDeson string `xml:"vICMSDeson,omitempty"` // 👈 ESSENCIAL
	VFCP       string `xml:"vFCP,omitempty"`
	VBCST      string `xml:"vBCST,omitempty"`
	VST        string `xml:"vST,omitempty"`

	VProd string `xml:"vProd"`

	VFrete string `xml:"vFrete,omitempty"`
	VSeg   string `xml:"vSeg,omitempty"`
	VDesc  string `xml:"vDesc,omitempty"`

	VII       string `xml:"vII,omitempty"`
	VIPI      string `xml:"vIPI,omitempty"`
	VIPIDevol string `xml:"vIPIDevol,omitempty"`

	VPIS    string `xml:"vPIS,omitempty"`
	VCOFINS string `xml:"vCOFINS,omitempty"`
	VOutro  string `xml:"vOutro,omitempty"`

	VNF string `xml:"vNF"`
}

type Transp struct {
	ModFrete string `xml:"modFrete"`
}

type Pag struct {
	DetPag []DetPag `xml:"detPag"`
}

type DetPag struct {
	TPag string `xml:"tPag"`
	VPag string `xml:"vPag"`
}

type PIS struct {
	PISOutr PISOutr `xml:"PISOutr"`
}

type PISOutr struct {
	CST  string `xml:"CST"`
	VBC  string `xml:"vBC"`
	PPIS string `xml:"pPIS"`
}

type COFINS struct {
	COFINSOutr COFINSOutr `xml:"COFINSOutr"`
}

type COFINSOutr struct {
	CST     string `xml:"CST"`
	VBC     string `xml:"vBC"`
	PCOFINS string `xml:"pCOFINS"`
}
