package validator

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	placaAntigaRegex   = regexp.MustCompile(`^[A-Z]{3}-?[0-9]{4}$`)
	placaMercosulRegex = regexp.MustCompile(`^[A-Z]{3}[0-9][A-Z][0-9]{2}$`)
)

// ValidateCPF valida um CPF usando o algoritmo oficial de dígitos verificadores.
func ValidateCPF(cpf string) bool {
	cpf = sanitizeDigits(cpf)
	if len(cpf) != 11 {
		return false
	}
	if allSameDigits(cpf) {
		return false
	}

	d1 := calcCPFDigit(cpf[:9], 10)
	d2 := calcCPFDigit(cpf[:9]+string(rune('0'+d1)), 11)
	return cpf[9] == byte('0'+d1) && cpf[10] == byte('0'+d2)
}

func calcCPFDigit(num string, weightStart int) int {
	sum := 0
	for i, c := range num {
		d, _ := strconv.Atoi(string(c))
		sum += d * (weightStart - i)
	}
	rest := sum % 11
	if rest < 2 {
		return 0
	}
	return 11 - rest
}

// ValidateCNPJ valida um CNPJ usando o algoritmo oficial de dígitos verificadores.
func ValidateCNPJ(cnpj string) bool {
	cnpj = sanitizeDigits(cnpj)
	if len(cnpj) != 14 {
		return false
	}
	if allSameDigits(cnpj) {
		return false
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	d1 := calcCNPJDigit(cnpj[:12], weights1)
	d2 := calcCNPJDigit(cnpj[:12]+string(rune('0'+d1)), weights2)
	return cnpj[12] == byte('0'+d1) && cnpj[13] == byte('0'+d2)
}

func calcCNPJDigit(num string, weights []int) int {
	sum := 0
	for i, c := range num {
		d, _ := strconv.Atoi(string(c))
		sum += d * weights[i]
	}
	rest := sum % 11
	if rest < 2 {
		return 0
	}
	return 11 - rest
}

// ValidatePlaca valida placas nos formatos brasileiro antigo (ABC-1234) e Mercosul (ABC1D23).
func ValidatePlaca(placa string) bool {
	placa = strings.ToUpper(strings.TrimSpace(placa))
	placa = strings.ReplaceAll(placa, "-", "")
	if placaAntigaRegex.MatchString(placa) {
		return true
	}
	return placaMercosulRegex.MatchString(placa)
}

// FormatPlaca normaliza a placa removendo hífen e convertendo para maiúsculas.
func FormatPlaca(placa string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(placa), "-", ""))
}

// SanitizeDocumento remove caracteres não numéricos de CPF/CNPJ antes da validação.
func SanitizeDocumento(doc string) string {
	return sanitizeDigits(doc)
}

// SanitizePlaca normaliza e remove caracteres inválidos de placas veiculares.
func SanitizePlaca(placa string) string {
	placa = strings.ToUpper(strings.TrimSpace(placa))
	var b strings.Builder
	for _, c := range placa {
		if (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func sanitizeDigits(s string) string {
	var b strings.Builder
	for _, c := range s {
		if c >= '0' && c <= '9' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func allSameDigits(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}
