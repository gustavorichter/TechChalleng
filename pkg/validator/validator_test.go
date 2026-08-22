package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCPF_Valido(t *testing.T) {
	assert.True(t, ValidateCPF("529.982.247-25"))
	assert.True(t, ValidateCPF("52998224725"))
}

func TestValidateCPF_Invalido(t *testing.T) {
	assert.False(t, ValidateCPF("111.111.111-11"))
	assert.False(t, ValidateCPF("12345678901"))
	assert.False(t, ValidateCPF("123"))
}

func TestValidateCNPJ_Valido(t *testing.T) {
	assert.True(t, ValidateCNPJ("04.252.011/0001-10"))
	assert.True(t, ValidateCNPJ("04252011000110"))
}

func TestValidateCNPJ_Invalido(t *testing.T) {
	assert.False(t, ValidateCNPJ("11.111.111/1111-11"))
	assert.False(t, ValidateCNPJ("123"))
}

func TestValidatePlaca(t *testing.T) {
	assert.True(t, ValidatePlaca("ABC-1234"))
	assert.True(t, ValidatePlaca("ABC1234"))
	assert.True(t, ValidatePlaca("ABC1D23"))
	assert.False(t, ValidatePlaca("INVALID"))
}
