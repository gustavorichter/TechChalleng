package valobj

import (
	"errors"
	"regexp"
	"strings"

	"github.com/techchalleng/oficina/pkg/validator"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// CPF representa um CPF validado.
type CPF struct {
	value string
}

func NewCPF(raw string) (CPF, error) {
	digits := sanitize(raw)
	if !validator.ValidateCPF(digits) {
		return CPF{}, errors.New("CPF inválido")
	}
	return CPF{value: digits}, nil
}

func (c CPF) String() string { return c.value }

// CNPJ representa um CNPJ validado.
type CNPJ struct {
	value string
}

func NewCNPJ(raw string) (CNPJ, error) {
	digits := sanitize(raw)
	if !validator.ValidateCNPJ(digits) {
		return CNPJ{}, errors.New("CNPJ inválido")
	}
	return CNPJ{value: digits}, nil
}

func (c CNPJ) String() string { return c.value }

// Placa representa uma placa de veículo validada.
type Placa struct {
	value string
}

func NewPlaca(raw string) (Placa, error) {
	if !validator.ValidatePlaca(raw) {
		return Placa{}, errors.New("placa inválida")
	}
	return Placa{value: validator.FormatPlaca(raw)}, nil
}

func (p Placa) String() string { return p.value }

// Email representa um endereço de e-mail validado.
type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if !emailRegex.MatchString(raw) {
		return Email{}, errors.New("e-mail inválido")
	}
	return Email{value: raw}, nil
}

func (e Email) String() string { return e.value }

func sanitize(s string) string {
	var b strings.Builder
	for _, c := range s {
		if c >= '0' && c <= '9' {
			b.WriteRune(c)
		}
	}
	return b.String()
}
