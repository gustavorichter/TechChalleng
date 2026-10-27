package entity

import (
	"testing"
	"time"

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

func TestOrdemServico_AprovarOrcamento(t *testing.T) {
	os := novaOSAguardandoAprovacao(t)
	os.Servicos = []ItemServicoOS{{ServicoID: uuid.New(), Nome: "Troca de óleo", ValorMaoObra: decimal.NewFromInt(100)}}

	err := os.AprovarOrcamento()
	require.NoError(t, err)
	assert.Equal(t, valobj.StatusEmExecucao, os.Status)
	require.NotNil(t, os.OrcamentoAprovado)
	assert.True(t, *os.OrcamentoAprovado)
}

func TestOrdemServico_AprovarOrcamento_Vazio(t *testing.T) {
	os := novaOSAguardandoAprovacao(t)
	err := os.AprovarOrcamento()
	require.Error(t, err)
	assert.Equal(t, valobj.StatusAguardandoAprovacao, os.Status)
}

func TestOrdemServico_RecusarOrcamento_VoltaAoDiagnostico(t *testing.T) {
	os := novaOSAguardandoAprovacao(t)
	err := os.RecusarOrcamento("  valor acima do combinado  ")
	require.NoError(t, err)
	assert.Equal(t, valobj.StatusEmDiagnostico, os.Status)
	assert.Equal(t, "valor acima do combinado", os.MotivoRecusa)
	require.NotNil(t, os.OrcamentoAprovado)
	assert.False(t, *os.OrcamentoAprovado)
}

func TestOrdemServico_RecusarOrcamento_SemMotivo(t *testing.T) {
	os := novaOSAguardandoAprovacao(t)
	err := os.RecusarOrcamento("   ")
	require.Error(t, err)
	assert.Equal(t, valobj.StatusAguardandoAprovacao, os.Status)
}

func TestOrdemServico_Finalizar_ExclusaoLogica(t *testing.T) {
	os := NewOrdemServico(uuid.New(), uuid.New(), "")
	for _, status := range []valobj.StatusOS{
		valobj.StatusEmDiagnostico,
		valobj.StatusAguardandoAprovacao,
		valobj.StatusEmExecucao,
		valobj.StatusFinalizada,
	} {
		require.NoError(t, os.AtualizarStatus(status))
	}
	require.NotNil(t, os.ExcluidaEm)
	assert.False(t, os.VisivelNaListagem())

	marcada := *os.ExcluidaEm
	require.NoError(t, os.AtualizarStatus(valobj.StatusEntregue))
	assert.Equal(t, marcada, *os.ExcluidaEm)
	assert.NotNil(t, os.EntregueEm)
}

func TestFiltrarEOrdenarListagem(t *testing.T) {
	base := time.Date(2026, 3, 1, 8, 0, 0, 0, time.UTC)
	execucaoNova := osComStatus(valobj.StatusEmExecucao, base.Add(3*time.Hour))
	execucaoAntiga := osComStatus(valobj.StatusEmExecucao, base.Add(1*time.Hour))
	aguardando := osComStatus(valobj.StatusAguardandoAprovacao, base.Add(4*time.Hour))
	diagnostico := osComStatus(valobj.StatusEmDiagnostico, base.Add(2*time.Hour))
	recebidaNova := osComStatus(valobj.StatusRecebida, base.Add(5*time.Hour))
	recebidaAntiga := osComStatus(valobj.StatusRecebida, base)
	finalizada := osComStatus(valobj.StatusFinalizada, base.Add(-time.Hour))
	entregue := osComStatus(valobj.StatusEntregue, base.Add(-2*time.Hour))
	now := time.Now().UTC()
	finalizada.ExcluidaEm = &now

	ordenadas := FiltrarEOrdenarListagem([]*OrdemServico{
		recebidaNova,
		finalizada,
		execucaoNova,
		diagnostico,
		entregue,
		aguardando,
		recebidaAntiga,
		execucaoAntiga,
	})

	require.Len(t, ordenadas, 6)
	assert.Equal(t, []*OrdemServico{
		execucaoAntiga,
		execucaoNova,
		aguardando,
		diagnostico,
		recebidaAntiga,
		recebidaNova,
	}, ordenadas)
}

func novaOSAguardandoAprovacao(t *testing.T) *OrdemServico {
	t.Helper()
	os := NewOrdemServico(uuid.New(), uuid.New(), "")
	require.NoError(t, os.AtualizarStatus(valobj.StatusEmDiagnostico))
	require.NoError(t, os.AtualizarStatus(valobj.StatusAguardandoAprovacao))
	return os
}

func osComStatus(status valobj.StatusOS, criadoEm time.Time) *OrdemServico {
	os := NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Status = status
	os.CriadoEm = criadoEm
	return os
}

func TestPeca_BaixarEstoque(t *testing.T) {
	peca := NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 10)
	err := peca.BaixarEstoque(3)
	require.NoError(t, err)
	assert.Equal(t, 7, peca.QuantidadeEstoque)

	err = peca.BaixarEstoque(10)
	require.Error(t, err)
}
