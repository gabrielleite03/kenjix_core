package nfe

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
)

type Builder struct {
	buf *bytes.Buffer
	enc *xml.Encoder
}

func NewBuilder() *Builder {
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	return &Builder{buf: buf, enc: enc}
}

func (b *Builder) Bytes() []byte {
	b.enc.Flush()
	return b.buf.Bytes()
}

func (b *Builder) start(name string, attrs ...xml.Attr) {
	b.enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: name},
		Attr: attrs,
	})
}

func (b *Builder) end(name string) {
	b.enc.EncodeToken(xml.EndElement{
		Name: xml.Name{Local: name},
	})
}

func (b *Builder) elem(name, value string) {
	if value == "" {
		return
	}
	b.enc.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: name},
	})
}

func (b *Builder) elemEmpty(name string) {
	start := xml.StartElement{Name: xml.Name{Local: name}}
	b.enc.EncodeToken(start)
	b.enc.EncodeToken(start.End())
}

// =========================
// BUILD PRINCIPAL
// =========================
func (b *Builder) BuildNFe(data NFeData) []byte {
	/*
		b.start("enviNFe",
			xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: "http://www.portalfiscal.inf.br/nfe"},
			xml.Attr{Name: xml.Name{Local: "versao"}, Value: "4.00"},
		)
		b.elem("idLote", data.IdLote)
		b.elem("indSinc", data.IndSinc)
	*/
	b.start("NFe",
		xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: "http://www.portalfiscal.inf.br/nfe"})

	// ⚠️ ID obrigatório (chave completa com "NFe")
	b.start("infNFe",
		xml.Attr{Name: xml.Name{Local: "Id"}, Value: data.ID},
		xml.Attr{Name: xml.Name{Local: "versao"}, Value: "4.00"},
	)

	//b.elem("idLote", data.IdLote)
	//b.elem("indSinc", data.IndSinc)

	// =========================
	// IDE (100% na ordem do XSD)
	// =========================
	b.start("ide")

	b.elem("cUF", UFToCUF(data.UF))
	b.elem("cNF", data.CNF)
	b.elem("natOp", data.NatOp)
	b.elem("mod", "55")
	b.elem("serie", data.Serie)
	b.elem("nNF", data.NNF)
	b.elem("dhEmi", data.DhEmi.Format("2006-01-02T15:04:05-07:00"))

	b.elem("tpNF", data.TpNF)
	b.elem("idDest", data.IdDest)
	b.elem("cMunFG", data.CMunFG)

	b.elem("tpImp", data.TpImp)
	b.elem("tpEmis", data.TpEmis)
	b.elem("cDV", strconv.Itoa(data.CDV))
	b.elem("tpAmb", data.TpAmb)

	b.elem("finNFe", data.FinNFe)
	b.elem("indFinal", data.IndFinal)
	b.elem("indPres", data.IndPres)

	//se for 1 adicionar informações do imtermediador
	b.elem("indIntermed", "0") // 0 = sem market place/ 1 ==com intermediario

	b.elem("procEmi", data.ProcEmi)
	b.elem("verProc", data.VerProc)

	b.end("ide")

	// =========================
	// EMIT
	// =========================
	b.start("emit")

	b.elem("CNPJ", data.EmitCNPJ)
	b.elem("CPF", data.EmitCPF)
	b.elem("xNome", data.EmitNome)

	b.start("enderEmit")
	b.elem("xLgr", data.EmitEnder.Logradouro)
	b.elem("nro", data.EmitEnder.Numero)
	b.elem("xBairro", data.EmitEnder.Bairro)
	b.elem("cMun", data.EmitEnder.CodigoMun)
	b.elem("xMun", data.EmitEnder.Municipio)
	b.elem("UF", data.EmitEnder.UF)
	b.elem("CEP", data.EmitEnder.CEP)
	b.elem("cPais", data.EmitEnder.CodigoPais)
	b.elem("xPais", data.EmitEnder.Pais)
	b.end("enderEmit")

	b.elem("IE", data.EmitIE)
	b.elem("CRT", data.EmitCRT) // Código de Regime Tributário 1=Simples Nacional;

	b.end("emit")

	// =========================
	// DEST
	// =========================
	b.start("dest")
	// se for CPF remover CNPJ
	//b.elem("CNPJ", data.DestCNPJ)
	b.elem("CPF", data.DestCPF)
	b.elem("xNome", data.DestNome)

	b.start("enderDest")
	b.elem("xLgr", data.DestEnder.Logradouro)
	b.elem("nro", data.DestEnder.Numero)
	b.elem("xBairro", data.DestEnder.Bairro)
	b.elem("cMun", data.DestEnder.CodigoMun)
	b.elem("xMun", data.DestEnder.Municipio)
	b.elem("UF", data.DestEnder.UF)
	b.elem("CEP", data.DestEnder.CEP)
	b.elem("cPais", data.DestEnder.CodigoPais)
	b.elem("xPais", data.DestEnder.Pais)
	b.end("enderDest")
	// se for CPF remover
	if data.DestCPF != "" {
		b.elem("indIEDest", data.DestIndIEDest)
	}

	// se for para CPF deve remover este campo
	if data.DestCNPJ != "" {
		b.elem("indIEDest", data.DestIndIEDest)
		b.elem("IE", data.EmitIE)
	}

	b.end("dest")

	// =========================
	// DET
	// =========================
	total := 0.0

	for i, item := range data.Itens {

		vProd := item.Qtd * item.Valor
		total += vProd

		b.start("det",
			xml.Attr{Name: xml.Name{Local: "nItem"}, Value: fmt.Sprintf("%d", i+1)},
		)

		b.start("prod")
		b.elem("cProd", item.Codigo)
		b.elem("cEAN", item.CEAN)
		b.elem("xProd", item.Desc)
		b.elem("NCM", item.NCM)
		//b.elem("CEST", "0100100") // koto qual o valor real?
		//b.elem("indEscala", "S")  // koto qual o valor real?
		b.elem("CFOP", item.CFOP)
		b.elem("uCom", item.Unidade)
		b.elem("qCom", fmt.Sprintf("%.4f", item.Qtd))
		b.elem("vUnCom", fmt.Sprintf("%.2f", item.Valor))
		b.elem("vProd", fmt.Sprintf("%.2f", vProd))
		b.elem("cEANTrib", item.CEANTrib)

		b.elem("uTrib", "UN")                              // koto qual o valor real?
		b.elem("qTrib", "1.0000")                          // koto qual o valor real?
		b.elem("vUnTrib", fmt.Sprintf("%.2f", item.Valor)) // koto qual o valor real?
		b.elem("indTot", "1")                              // koto qual o valor real?

		b.end("prod")

		b.start("imposto")

		// ICMS SN
		b.start("ICMS")
		b.start("ICMSSN102")
		b.elem("orig", "0")
		b.elem("CSOSN", "102")
		b.end("ICMSSN102")
		b.end("ICMS")

		// PIS
		b.start("PIS")
		/*
			b.start("PISOutr")
			b.elem("CST", "99")
			b.elem("vBC", "0.00")
			b.elem("pPIS", "0.00")
			b.elem("vPIS", "0.00")
			b.end("PISOutr")
		*/
		b.start("PISNT")
		b.elem("CST", "07")
		b.end("PISNT")

		/*
			b.start("PISQtde")
			b.elem("vPIS", "0.00")
			b.end("PISQtde")
			b.elem("CST", "03")
			b.elem("qBCProd", "1")
			b.elem("vAliqProd", "0.00")
			b.elem("vPIS", "0.00")
			b.start("PISQtde")
			b.elem("CST", "04")
			b.end("PISQtde")
		*/
		b.end("PIS")

		// COFINS
		b.start("COFINS")
		/*
			b.start("COFINSOutr")
			b.elem("CST", "99")
			b.elem("vBC", "0.00")
			b.elem("pCOFINS", "0.00")
			b.elem("vCOFINS", "0.00")
			b.end("COFINSOutr")
		*/
		b.start("COFINSNT")
		b.elem("CST", "07")
		b.end("COFINSNT")

		b.end("COFINS")

		b.end("imposto")
		b.end("det")
	}

	// =========================
	// TOTAL
	// =========================
	b.start("total")
	b.start("ICMSTot")

	b.elem("vBC", "0.00")
	b.elem("vICMS", "0.00")
	b.elem("vICMSDeson", "0.00")

	b.elem("vFCP", "0.00")
	b.elem("vBCST", "0.00")
	b.elem("vST", "0.00")
	b.elem("vFCPST", "0.00")
	b.elem("vFCPSTRet", "0.00")
	b.elem("vProd", fmt.Sprintf("%.2f", total))

	b.elem("vFrete", "0.00")
	b.elem("vSeg", "0.00")
	b.elem("vDesc", "0.00")
	b.elem("vII", "0.00")

	b.elem("vIPI", "0.00")
	b.elem("vIPIDevol", "0.00")
	b.elem("vPIS", "0.00")
	b.elem("vCOFINS", "0.00")
	b.elem("vOutro", "0.00")
	b.elem("vNF", fmt.Sprintf("%.2f", total))

	b.end("ICMSTot")
	b.end("total")

	// =========================
	// TRANSP
	// =========================
	b.start("transp")
	b.elem("modFrete", "9")
	b.end("transp")

	// =========================
	// COBR
	// =========================
	/*
		b.start("cobr")
		b.start("fat")
		b.elem("nFat", "123")
		b.elem("vOrig", fmt.Sprintf("%.2f", total))
		b.elem("vLiq", fmt.Sprintf("%.2f", total))
		b.end("fat")
		b.end("cobr")
	*/
	// =========================
	// PAG
	// =========================
	b.start("pag")

	for _, p := range data.Pagamentos {
		b.start("detPag")
		b.elem("tPag", p.Tipo)
		b.elem("vPag", fmt.Sprintf("%.2f", p.Valor))
		b.end("detPag")
		b.elem("vTroco", "0.00")
	}

	b.end("pag")

	b.end("infNFe")
	b.start("Signature", xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: "http://www.w3.org/2000/09/xmldsig#"})
	b.start("SignedInfo")
	b.start("CanonicalizationMethod", xml.Attr{Name: xml.Name{Local: "Algorithm"}, Value: "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"})
	b.end("CanonicalizationMethod")
	b.start("SignatureMethod", xml.Attr{Name: xml.Name{Local: "Algorithm"}, Value: "http://www.w3.org/2000/09/xmldsig#rsa-sha1"})
	b.end("SignatureMethod")
	b.start("Reference", xml.Attr{Name: xml.Name{Local: "URI"}, Value: "#" + data.ID})
	b.start("Transforms")
	b.start("Transform", xml.Attr{Name: xml.Name{Local: "Algorithm"}, Value: "http://www.w3.org/2000/09/xmldsig#enveloped-signature"})
	b.end("Transform")
	b.start("Transform", xml.Attr{Name: xml.Name{Local: "Algorithm"}, Value: "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"})
	b.end("Transform")
	b.end("Transforms")
	b.start("DigestMethod", xml.Attr{Name: xml.Name{Local: "Algorithm"}, Value: "http://www.w3.org/2000/09/xmldsig#sha1"})
	b.end("DigestMethod")
	b.start("DigestValue")
	b.end("DigestValue")
	b.end("Reference")
	b.end("SignedInfo")
	b.start("SignatureValue")
	b.end("SignatureValue")
	b.start("KeyInfo")
	b.start("X509Data")
	b.start("X509Certificate")
	b.end("X509Certificate")
	b.end("X509Data")
	b.end("KeyInfo")
	b.end("Signature")
	b.end("NFe")

	return b.Bytes()
}

func ValidateNFe(data NFeData) error {
	var errs ValidationErrors

	validateIde(data, &errs)
	validateItens(data.Itens, &errs)
	validatePagamentos(data.Pagamentos, &errs)

	if len(errs) > 0 {
		return errs
	}

	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var msg string
	for _, e := range v {
		msg += e.Field + ": " + e.Message + "\n"
	}
	return msg
}

func validatePagamentos(pags []Pagamento, errs *ValidationErrors) {

	if len(pags) == 0 {
		*errs = append(*errs, ValidationError{"pag", "deve conter pelo menos 1 pagamento"})
		return
	}

	totalPag := 0.0

	for i, p := range pags {

		prefix := fmt.Sprintf("pag[%d]", i+1)

		if p.Tipo == "" {
			*errs = append(*errs, ValidationError{prefix + ".tPag", "obrigatório"})
		}

		if p.Valor <= 0 {
			*errs = append(*errs, ValidationError{prefix + ".vPag", "deve ser maior que 0"})
		}

		totalPag += p.Valor
	}

	// (Opcional) validar soma depois se quiser cruzar com total da nota
	if totalPag <= 0 {
		*errs = append(*errs, ValidationError{"pag.total", "valor total inválido"})
	}
}

func validateItens(itens []Item, errs *ValidationErrors) {

	if len(itens) == 0 {
		*errs = append(*errs, ValidationError{"det", "deve conter pelo menos 1 item"})
		return
	}

	for i, item := range itens {

		prefix := fmt.Sprintf("det[%d]", i+1)

		if item.Codigo == "" {
			*errs = append(*errs, ValidationError{prefix + ".cProd", "obrigatório"})
		}

		if item.Desc == "" {
			*errs = append(*errs, ValidationError{prefix + ".xProd", "obrigatório"})
		}

		if len(item.NCM) != 8 {
			*errs = append(*errs, ValidationError{prefix + ".NCM", "deve ter 8 dígitos"})
		}

		if item.CFOP == "" {
			*errs = append(*errs, ValidationError{prefix + ".CFOP", "obrigatório"})
		}

		if item.Unidade == "" {
			*errs = append(*errs, ValidationError{prefix + ".uCom", "obrigatório"})
		}

		if item.Qtd <= 0 {
			*errs = append(*errs, ValidationError{prefix + ".qCom", "deve ser maior que 0"})
		}

		if item.Valor <= 0 {
			*errs = append(*errs, ValidationError{prefix + ".vUnCom", "deve ser maior que 0"})
		}
	}
}

func validateIde(data NFeData, errs *ValidationErrors) {

	if data.UF == "" {
		*errs = append(*errs, ValidationError{"ide.cUF", "obrigatório"})
	}

	if data.Serie == "" {
		*errs = append(*errs, ValidationError{"ide.serie", "obrigatório"})
	}

	if len(data.NNF) < 1 {
		*errs = append(*errs, ValidationError{"ide.nNF", "deve ter 9 dígitos"})
	}

	if data.TpAmb != "1" && data.TpAmb != "2" {
		*errs = append(*errs, ValidationError{"ide.tpAmb", "valores permitidos: 1,2"})
	}

	if data.TpEmis == "" {
		*errs = append(*errs, ValidationError{"ide.tpEmis", "obrigatório"})
	}

	if data.NatOp == "" {
		*errs = append(*errs, ValidationError{"ide.natOp", "obrigatório"})
	}

	if data.CMunFG == "" {
		*errs = append(*errs, ValidationError{"ide.cMunFG", "obrigatório"})
	}
}

func validateProd(p ProdData, errs *ValidationErrors) {

	if p.CProd == "" {
		*errs = append(*errs, ValidationError{"prod.cProd", "obrigatório"})
	}

	if p.XProd == "" {
		*errs = append(*errs, ValidationError{"prod.xProd", "obrigatório"})
	}

	if len(p.NCM) != 8 {
		*errs = append(*errs, ValidationError{"prod.NCM", "deve ter 8 dígitos"})
	}

	if p.VProd == "" {
		*errs = append(*errs, ValidationError{"prod.vProd", "obrigatório"})
	}
}

var pisAliqRegex = regexp.MustCompile(`^0(\.\d{2,4})?$|^[1-9]\d{0,2}(\.\d{2,4})?$`)

func validatePIS(p PISData, errs *ValidationErrors) {

	if p.CST == "" {
		*errs = append(*errs, ValidationError{"PIS.CST", "obrigatório"})
	}

	if p.VBC == "" {
		*errs = append(*errs, ValidationError{"PIS.vBC", "obrigatório"})
	}

	if !pisAliqRegex.MatchString(p.PPIS) {
		*errs = append(*errs, ValidationError{"PIS.pPIS", "formato inválido"})
	}
}

func validateCOFINS(c COFINSData, errs *ValidationErrors) {

	if c.CST == "" {
		*errs = append(*errs, ValidationError{"COFINS.CST", "obrigatório"})
	}

	if c.VBC == "" {
		*errs = append(*errs, ValidationError{"COFINS.vBC", "obrigatório"})
	}

	if !pisAliqRegex.MatchString(c.PCOFINS) {
		*errs = append(*errs, ValidationError{"COFINS.pCOFINS", "formato inválido"})
	}
}
func validateTotal(t TotData, errs *ValidationErrors) {

	if t.VNF == "" {
		*errs = append(*errs, ValidationError{"total.vNF", "obrigatório"})
	}

	if t.VProd == "" {
		*errs = append(*errs, ValidationError{"total.vProd", "obrigatório"})
	}
}
