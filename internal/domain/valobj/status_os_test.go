package valobj

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusOS_TransicoesValidas(t *testing.T) {
	tests := []struct {
		atual   StatusOS
		destino StatusOS
		valido  bool
	}{
		{StatusRecebida, StatusEmDiagnostico, true},
		{StatusEmDiagnostico, StatusAguardandoAprovacao, true},
		{StatusAguardandoAprovacao, StatusEmExecucao, true},
		{StatusEmExecucao, StatusFinalizada, true},
		{StatusFinalizada, StatusEntregue, true},
		{StatusRecebida, StatusFinalizada, false},
		{StatusEntregue, StatusRecebida, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.atual)+"_para_"+string(tt.destino), func(t *testing.T) {
			_, err := tt.atual.TransicionarPara(tt.destino)
			if tt.valido {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestStatusOS_ProximoStatus(t *testing.T) {
	proximo, ok := StatusRecebida.ProximoStatus()
	assert.True(t, ok)
	assert.Equal(t, StatusEmDiagnostico, proximo)

	_, ok = StatusEntregue.ProximoStatus()
	assert.False(t, ok)
}

func TestNewStatusOS_Invalido(t *testing.T) {
	_, err := NewStatusOS("Status Inexistente")
	assert.Error(t, err)
}
