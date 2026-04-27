package nfe

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Gera a chave de acesso da NF-e (44 dígitos)
func GenerateNFeKey(uf string, emissionDate time.Time, cnpj string, model string, serie string, nNF string, tpEmissao string, cNF string) (string, string, int) {
	cUF := UFToCUF(uf)                                   // 2 dígitos
	aa := fmt.Sprintf("%02d", emissionDate.Year()%100)   // 2 dígitos
	mm := fmt.Sprintf("%02d", int(emissionDate.Month())) // 2 dígitos
	cnpj14 := padLeft(onlyDigits(cnpj), 14)              // 14 dígitos

	// Usando padLeft para garantir os zeros à esquerda corretamente
	mod := padLeft(onlyDigits(model), 2)    // 2 dígitos
	s := padLeft(onlyDigits(serie), 3)      // 3 dígitos
	n := padLeft(onlyDigits(nNF), 9)        // 9 dígitos
	tp := padLeft(onlyDigits(tpEmissao), 1) // 1 dígito
	c := padLeft(onlyDigits(cNF), 8)        // 8 dígitos

	// Montagem: 2+2+2+14+2+3+9+1+8 = 43 dígitos exatos
	chaveSemDV := cUF + aa + mm + cnpj14 + mod + s + n + tp + c

	// DEBUG (Opcional): fmt.Println("LEN:", len(chaveSemDV))
	// Deve ser sempre 43

	dv := calculateDV(chaveSemDV)

	return chaveSemDV + strconv.Itoa(dv), c, dv
}

func calculateDV(chave string) int {
	// 1. Limpeza: Garante que estamos lidando apenas com os dígitos
	// Remove o prefixo "NFe" caso tenha sido passado por engano
	chave = strings.ReplaceAll(chave, "NFe", "")

	// Remove qualquer caractere que não seja número (espaços, traços, etc)
	reg := regexp.MustCompile(`[^0-9]`)
	chave = reg.ReplaceAllString(chave, "")

	// A chave DEVE ter 43 dígitos para calcular o 44º
	if len(chave) != 43 {
		// Se chegar aqui com tamanho errado, o cálculo falhará silenciosamente
		// ou retornará 0, o que causa a Rejeição 225.
		return 0
	}

	soma := 0
	peso := 2

	// Cálculo Módulo 11 (Peso de 2 a 9, da direita para a esquerda)
	for i := len(chave) - 1; i >= 0; i-- {
		num := int(chave[i] - '0') // Conversão mais eficiente que Atoi
		soma += num * peso
		peso++
		if peso > 9 {
			peso = 2
		}
	}

	resto := soma % 11

	// Regra da SEFAZ para NF-e
	if resto < 2 {
		return 0
	}
	return 11 - resto
}

func padLeft(value string, totalLen int) string {
	if len(value) >= totalLen {
		return value[0:totalLen] // Corta se for maior
	}
	return strings.Repeat("0", totalLen-len(value)) + value
}

func onlyDigits(s string) string {
	out := ""
	for _, r := range s {
		if r >= '0' && r <= '9' {
			out += string(r)
		}
	}
	return out
}

func UFToCUF(uf string) string {
	switch uf {
	case "RO":
		return "11"
	case "AC":
		return "12"
	case "AM":
		return "13"
	case "RR":
		return "14"
	case "PA":
		return "15"
	case "AP":
		return "16"
	case "TO":
		return "17"

	case "MA":
		return "21"
	case "PI":
		return "22"
	case "CE":
		return "23"
	case "RN":
		return "24"
	case "PB":
		return "25"
	case "PE":
		return "26"
	case "AL":
		return "27"
	case "SE":
		return "28"
	case "BA":
		return "29"

	case "MG":
		return "31"
	case "ES":
		return "32"
	case "RJ":
		return "33"
	case "SP":
		return "35"

	case "PR":
		return "41"
	case "SC":
		return "42"
	case "RS":
		return "43"

	case "MS":
		return "50"
	case "MT":
		return "51"
	case "GO":
		return "52"
	case "DF":
		return "53"

	default:
		return "35" // fallback (SP)
	}
}
