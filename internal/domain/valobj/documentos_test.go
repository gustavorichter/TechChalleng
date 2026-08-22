package valobj

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCPF_Valido(t *testing.T) {
	cpf, err := NewCPF("529.982.247-25")
	require.NoError(t, err)
	assert.Equal(t, "52998224725", cpf.String())
}

func TestNewCPF_Invalido(t *testing.T) {
	_, err := NewCPF("111.111.111-11")
	assert.Error(t, err)
}

func TestNewCNPJ_Valido(t *testing.T) {
	cnpj, err := NewCNPJ("04.252.011/0001-10")
	require.NoError(t, err)
	assert.Equal(t, "04252011000110", cnpj.String())
}

func TestNewPlaca(t *testing.T) {
	placa, err := NewPlaca("ABC-1234")
	require.NoError(t, err)
	assert.Equal(t, "ABC1234", placa.String())

	placa, err = NewPlaca("ABC1D23")
	require.NoError(t, err)
	assert.Equal(t, "ABC1D23", placa.String())
}

func TestNewEmail(t *testing.T) {
	email, err := NewEmail("  Test@Example.COM  ")
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", email.String())

	_, err = NewEmail("invalid")
	assert.Error(t, err)
}
