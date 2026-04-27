package nfe

import "time"

type IdeData struct {
	CUF, CNF, NatOp, Mod, Serie, NNF string
	DhEmi, TpNF, IdDest, CMunFG      string
	TpImp, TpEmis, CDV               string
	TpAmb, FinNFe, IndFinal, IndPres string
	ProcEmi, VerProc                 string
}

type ProdData struct {
	CProd, CEAN, XProd, NCM, CFOP string
	UCom, QCom, VUnCom, VProd     string
}

type PISData struct {
	CST, VBC, PPIS string
}

type COFINSData struct {
	CST, VBC, PCOFINS string
}

type TotData struct {
	VBC, VICMS, VICMSDeson, VProd, VNF string
}
type NFeData struct {
	// ===== CONTROLE =====
	IdLote  string
	IndSinc string

	// ===== CHAVE =====
	ID  string // NFe351...
	CNF string // 8 dígitos

	// ===== IDE =====
	UF          string
	NatOp       string
	Serie       string
	NNF         string
	DhEmi       time.Time
	TpNF        string
	IdDest      string
	CMunFG      string
	TpImp       string
	TpEmis      string
	CDV         int
	TpAmb       string
	FinNFe      string
	IndFinal    string
	IndPres     string
	IndIntermed string
	ProcEmi     string
	VerProc     string

	// ===== EMITENTE =====
	EmitCNPJ  string
	EmitCPF   string
	EmitNome  string
	EmitIE    string
	EmitCRT   string
	EmitEnder Endereco

	// ===== DEST =====
	DestCNPJ      string
	DestCPF       string
	DestNome      string
	DestIndIEDest string
	DestEnder     Endereco

	// ===== PRODUTOS =====
	Itens []Item

	// ===== PAGAMENTO =====
	Pagamentos []Pagamento
}

type Endereco struct {
	Logradouro string
	Numero     string
	Bairro     string
	CodigoMun  string
	Municipio  string
	UF         string
	CEP        string
	CodigoPais string
	Pais       string
}

type Item struct {
	Codigo   string
	CEAN     string
	Desc     string
	NCM      string
	CFOP     string
	Unidade  string
	Qtd      float64
	Valor    float64
	CEANTrib string
}

type Pagamento struct {
	Tipo  string
	Valor float64
}
