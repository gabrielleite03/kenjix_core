package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gabrielleite03/kenjix_core/cmd/api/router"
	"github.com/gabrielleite03/kenjix_core/internal/service"
	"github.com/gabrielleite03/kenjix_core/internal/service/nfe"
)

func main() {
	r := router.NewRouter()
	r.Register()
	log.Println("Kenjix Core iniciado na porta 7010")
	err := http.ListenAndServe(":7010", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func mainNFe() {
	randomCNF := fmt.Sprintf("%08d", rand.Intn(99999999))
	now := time.Now()
	// =========================
	// 1. Gerar chave + cNF (igual você já faz)
	// =========================
	chave, cNF, dv := nfe.GenerateNFeKey(
		"SP",
		now,
		"65468523000102",
		"55",
		"1",
		"1",
		"1",
		randomCNF,
	)

	// ⚠️ obrigatório prefixo
	id := "NFe" + chave

	// =========================
	// 2. Montar NFeData (equivalente ao seu DTO)
	// =========================
	data := nfe.NFeData{
		IdLote:  "2",
		IndSinc: "1",

		ID:    id,
		CNF:   cNF,
		DhEmi: now,

		UF:          "SP",
		NatOp:       "Venda",
		Serie:       "1",
		NNF:         "1", // ⚠️ manter 9 dígitos
		TpNF:        "1",
		IdDest:      "1",
		CMunFG:      "3550308",
		TpImp:       "1",
		TpEmis:      "1",
		CDV:         dv,
		TpAmb:       "2",
		FinNFe:      "1",
		IndFinal:    "1",
		IndPres:     "1",
		IndIntermed: "1",
		ProcEmi:     "0",
		VerProc:     "1.0",

		EmitCNPJ: "65468523000102",
		//	EmitCPF:  "32843874807",
		EmitNome: "KENJI IMPORTACAO E COMERCIO LTDA",
		EmitIE:   "158447676112",
		EmitCRT:  "1",
		EmitEnder: nfe.Endereco{
			Logradouro: "Rua A",
			Numero:     "100",
			Bairro:     "Centro",
			CodigoMun:  "3550308",
			Municipio:  "Sao Paulo",
			UF:         "SP",
			CEP:        "01001000",
			CodigoPais: "1058",
			Pais:       "Brasil",
		},

		//	DestCNPJ:      "99999999000199",
		DestCPF:       "32843874807",
		DestNome:      "NF-E EMITIDA EM AMBIENTE DE HOMOLOGACAO - SEM VALOR FISCAL",
		DestIndIEDest: "9",
		DestEnder: nfe.Endereco{
			Logradouro: "Rua B",
			Numero:     "200",
			Bairro:     "Centro",
			CodigoMun:  "3550308",
			Municipio:  "Sao Paulo",
			UF:         "SP",
			CEP:        "01002000",
			CodigoPais: "1058",
			Pais:       "Brasil",
		},

		Itens: []nfe.Item{
			{
				Codigo:   "001",
				Desc:     "Produto Teste",
				CEAN:     "SEM GTIN",
				NCM:      "40169990",
				CFOP:     "5102",
				Unidade:  "UN",
				Qtd:      1,
				Valor:    10.00,
				CEANTrib: "SEM GTIN",
			},
		},

		Pagamentos: []nfe.Pagamento{
			{
				Tipo:  "01",
				Valor: 10.00,
			},
		},
	}

	// =========================
	// 3. (OPCIONAL) validar antes
	// =========================
	if err := nfe.ValidateNFe(data); err != nil {
		fmt.Println("ERRO DE VALIDAÇÃO:", err)
		return
	}

	// =========================
	// 4. Gerar XML com Builder
	// =========================
	builder := nfe.NewBuilder()
	xmlBytes := builder.BuildNFe(data)

	// =========================
	// 5. Resultado
	// =========================
	fmt.Println(string(xmlBytes))

	_, err := service.EmitirNFeKoto(xmlBytes)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			fmt.Printf("URL ERROR: %+v\n", urlErr)
			if opErr, ok := urlErr.Err.(*net.OpError); ok {
				fmt.Printf("NET ERROR: %+v\n", opErr)
			}
		}
		panic(err)
	}

	//fmt.Println(string(resp))
}
