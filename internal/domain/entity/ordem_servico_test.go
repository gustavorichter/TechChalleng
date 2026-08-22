package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

func TestOrdemServico_AtualizarStatus_Valido(t *testing.T) {
	os := NewOrdemServico(uuid.New(), uuid.New(), "teste")
	err := os.AtualizarStatus(valobj.StatusEmDiagnostico)
	require.NoError(t, err)
	assert.Equal(t, valobj.StatusEmDiagnostico, os.Status)
}

func TestOrdemServico_AtualizarStatus_Invalido(t *testing.T) {
	os := NewOrdemServico(uuid.New(), uuid.New(), "teste")
	err := os.AtualizarStatus(valobj.StatusFinalizada)
	require.Error(t, err)
	assert.Equal(t, valobj.StatusRecebida, os.Status)
}

func TestOrdemServico_CicloCompleto(t *testing.T) {
	os := NewOrdemServico(uuid.New(), uuid.New(), "teste")
	statuses := []valobj.StatusOS{
		valobj.StatusEmDiagnostico,
		valobj.StatusAguardandoAprovacao,
		valobj.StatusEmExecucao,
		valobj.StatusFinalizada,
		valobj.StatusEntregue,
	}
	for _, s := range statuses {
		err := os.AtualizarStatus(s)
		require.NoError(t, err, "falha ao transicionar para %s", s)
	}
	assert.NotNil(t, os.EntregueEm)
	_, ok := os.TempoExecucao()
	assert.True(t, ok)
}

func TestOrdemServico_AdicionarPeca_StatusInvalido(t *testing.T) {
	os := NewOrdemServico(uuid.New(), uuid.New(), "")
	_ = os.AtualizarStatus(valobj.StatusEmDiagnostico)
	_ = os.AtualizarStatus(valobj.StatusAguardandoAprovacao)

	err := os.AdicionarPeca(ItemPecaOS{
		PecaID: uuid.New(), Nome: "Filtro", Quantidade: 1,
		ValorUnitario: decimal.NewFromFloat(10),
	})
	require.Error(t, err)
}

func TestPeca_BaixarEstoque(t *testing.T) {
	peca := NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 10)
	err := peca.BaixarEstoque(3)
	require.NoError(t, err)
	assert.Equal(t, 7, peca.QuantidadeEstoque)

	err = peca.BaixarEstoque(10)
	require.Error(t, err)
}
