package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techchalleng/oficina/internal/application/dto"
	"github.com/techchalleng/oficina/internal/application/port"
	"github.com/techchalleng/oficina/internal/application/usecase"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/valobj"
	"github.com/techchalleng/oficina/internal/infra/http/handler"
)

type memOS struct {
	ordens map[uuid.UUID]*entity.OrdemServico
}

func (m *memOS) Criar(_ context.Context, os *entity.OrdemServico) error {
	m.ordens[os.ID] = os
	return nil
}
func (m *memOS) BuscarPorID(_ context.Context, id uuid.UUID) (*entity.OrdemServico, error) {
	os, ok := m.ordens[id]
	if !ok {
		return nil, assert.AnError
	}
	return os, nil
}
func (m *memOS) Listar(_ context.Context) ([]*entity.OrdemServico, error) {
	result := make([]*entity.OrdemServico, 0, len(m.ordens))
	for _, os := range m.ordens {
		result = append(result, os)
	}
	return result, nil
}
func (m *memOS) Atualizar(_ context.Context, os *entity.OrdemServico) error {
	m.ordens[os.ID] = os
	return nil
}
func (m *memOS) TempoMedioExecucao(context.Context) (float64, error) { return 2.5, nil }

type memClientes struct{ itens map[uuid.UUID]*entity.Cliente }

func (m *memClientes) Criar(context.Context, *entity.Cliente) error { return nil }
func (m *memClientes) BuscarPorID(_ context.Context, id uuid.UUID) (*entity.Cliente, error) {
	c, ok := m.itens[id]
	if !ok {
		return nil, assert.AnError
	}
	return c, nil
}
func (m *memClientes) Listar(context.Context) ([]*entity.Cliente, error) { return nil, nil }
func (m *memClientes) Atualizar(context.Context, *entity.Cliente) error  { return nil }
func (m *memClientes) Excluir(context.Context, uuid.UUID) error          { return nil }

type memVeiculos struct{ itens map[uuid.UUID]*entity.Veiculo }

func (m *memVeiculos) Criar(context.Context, *entity.Veiculo) error { return nil }
func (m *memVeiculos) BuscarPorID(_ context.Context, id uuid.UUID) (*entity.Veiculo, error) {
	v, ok := m.itens[id]
	if !ok {
		return nil, assert.AnError
	}
	return v, nil
}
func (m *memVeiculos) ListarPorCliente(context.Context, uuid.UUID) ([]*entity.Veiculo, error) {
	return nil, nil
}
func (m *memVeiculos) Listar(context.Context) ([]*entity.Veiculo, error) { return nil, nil }
func (m *memVeiculos) Atualizar(context.Context, *entity.Veiculo) error  { return nil }
func (m *memVeiculos) Excluir(context.Context, uuid.UUID) error          { return nil }

type memServicos struct{ itens map[uuid.UUID]*entity.Servico }

func (m *memServicos) Criar(context.Context, *entity.Servico) error { return nil }
func (m *memServicos) BuscarPorID(_ context.Context, id uuid.UUID) (*entity.Servico, error) {
	s, ok := m.itens[id]
	if !ok {
		return nil, assert.AnError
	}
	return s, nil
}
func (m *memServicos) Listar(context.Context) ([]*entity.Servico, error) { return nil, nil }
func (m *memServicos) Atualizar(context.Context, *entity.Servico) error  { return nil }
func (m *memServicos) Excluir(context.Context, uuid.UUID) error          { return nil }

type memPecas struct{ itens map[uuid.UUID]*entity.Peca }

func (m *memPecas) Criar(context.Context, *entity.Peca) error { return nil }
func (m *memPecas) BuscarPorID(_ context.Context, id uuid.UUID) (*entity.Peca, error) {
	p, ok := m.itens[id]
	if !ok {
		return nil, assert.AnError
	}
	return p, nil
}
func (m *memPecas) Listar(context.Context) ([]*entity.Peca, error) { return nil, nil }
func (m *memPecas) Atualizar(_ context.Context, p *entity.Peca) error {
	m.itens[p.ID] = p
	return nil
}
func (m *memPecas) Excluir(context.Context, uuid.UUID) error { return nil }

type memNotificador struct{ total int }

func (m *memNotificador) NotificarAtualizacaoStatus(context.Context, port.NotificacaoStatus) error {
	m.total++
	return nil
}

func TestIntegracaoHTTP_CicloDeVidaOS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	email, err := valobj.NewEmail("joao@email.com")
	require.NoError(t, err)
	cliente := entity.NewCliente("João Silva", email, "11999999999")
	veiculo := &entity.Veiculo{ID: uuid.New(), ClienteID: cliente.ID, Marca: "VW", Modelo: "Gol", Ano: 2022}
	servico := entity.NewServico("Troca de óleo", "Completa", decimal.NewFromInt(150))
	peca := entity.NewPeca("Filtro", "FLT001", decimal.NewFromInt(45), 10)

	osRepo := &memOS{ordens: map[uuid.UUID]*entity.OrdemServico{}}
	notificador := &memNotificador{}
	uc := usecase.NewOrdemServicoUseCase(
		osRepo,
		&memClientes{itens: map[uuid.UUID]*entity.Cliente{cliente.ID: cliente}},
		&memVeiculos{itens: map[uuid.UUID]*entity.Veiculo{veiculo.ID: veiculo}},
		&memServicos{itens: map[uuid.UUID]*entity.Servico{servico.ID: servico}},
		&memPecas{itens: map[uuid.UUID]*entity.Peca{peca.ID: peca}},
		notificador,
	)
	h := handler.NewOrdemServicoHandler(uc)

	r := gin.New()
	r.POST("/ordens-servico", h.Criar)
	r.GET("/ordens-servico", h.Listar)
	r.GET("/ordens-servico/:id/status", h.AcompanharStatus)
	r.PUT("/ordens-servico/:id/status", h.AtualizarStatus)
	r.POST("/ordens-servico/:id/orcamento", h.DecidirOrcamento)

	criada := postJSON[dto.OrdemServicoResponse](t, r, "/ordens-servico", dto.CriarOSRequest{
		ClienteID:   cliente.ID,
		VeiculoID:   veiculo.ID,
		Observacoes: "Barulho no motor",
		Servicos:    []dto.AdicionarServicoOSRequest{{ServicoID: servico.ID}},
		Pecas:       []dto.AdicionarPecaOSRequest{{PecaID: peca.ID, Quantidade: 2}},
	}, http.StatusCreated)
	assert.Equal(t, "Recebida", criada.Status)
	assert.True(t, criada.ValorTotal.Equal(decimal.NewFromInt(240)))

	osRepo.ordens[criada.ID].CriadoEm = time.Date(2026, 2, 1, 8, 0, 0, 0, time.UTC)

	segunda := postJSON[dto.OrdemServicoResponse](t, r, "/ordens-servico", dto.CriarOSRequest{
		ClienteID: cliente.ID,
		VeiculoID: veiculo.ID,
	}, http.StatusCreated)
	osRepo.ordens[segunda.ID].CriadoEm = time.Date(2026, 2, 2, 8, 0, 0, 0, time.UTC)

	putJSON[dto.OrdemServicoResponse](t, r, "/ordens-servico/"+criada.ID.String()+"/status", dto.AtualizarStatusOSRequest{Status: "Em diagnóstico"}, http.StatusOK)
	putJSON[dto.OrdemServicoResponse](t, r, "/ordens-servico/"+criada.ID.String()+"/status", dto.AtualizarStatusOSRequest{Status: "Aguardando aprovação"}, http.StatusOK)

	status := getJSON[dto.StatusOSResponse](t, r, "/ordens-servico/"+criada.ID.String()+"/status", http.StatusOK)
	assert.Equal(t, "Aguardando aprovação", status.Status)
	assert.True(t, status.ValorTotal.Equal(decimal.NewFromInt(240)))

	recusada := postJSON[dto.OrdemServicoResponse](t, r, "/ordens-servico/"+criada.ID.String()+"/orcamento", dto.DecidirOrcamentoRequest{
		Aprovado: false,
		Motivo:   "valor acima do combinado",
	}, http.StatusOK)
	assert.Equal(t, "Em diagnóstico", recusada.Status)

	putJSON[dto.OrdemServicoResponse](t, r, "/ordens-servico/"+criada.ID.String()+"/status", dto.AtualizarStatusOSRequest{Status: "Aguardando aprovação"}, http.StatusOK)
	aprovada := postJSON[dto.OrdemServicoResponse](t, r, "/ordens-servico/"+criada.ID.String()+"/orcamento", dto.DecidirOrcamentoRequest{Aprovado: true}, http.StatusOK)
	assert.Equal(t, "Em execução", aprovada.Status)

	fila := getJSON[[]dto.OrdemServicoResponse](t, r, "/ordens-servico", http.StatusOK)
	require.Len(t, fila, 2)
	assert.Equal(t, criada.ID, fila[0].ID)
	assert.Equal(t, segunda.ID, fila[1].ID)

	finalizada := putJSON[dto.OrdemServicoResponse](t, r, "/ordens-servico/"+criada.ID.String()+"/status", dto.AtualizarStatusOSRequest{Status: "Finalizada"}, http.StatusOK)
	assert.Equal(t, "Finalizada", finalizada.Status)
	require.NotNil(t, finalizada.ExcluidaEm)
	assert.Equal(t, 8, peca.QuantidadeEstoque)

	fila = getJSON[[]dto.OrdemServicoResponse](t, r, "/ordens-servico", http.StatusOK)
	require.Len(t, fila, 1)
	assert.Equal(t, segunda.ID, fila[0].ID)

	publica := getJSON[dto.StatusOSResponse](t, r, "/ordens-servico/"+criada.ID.String()+"/status", http.StatusOK)
	assert.Equal(t, "Finalizada", publica.Status)
	assert.GreaterOrEqual(t, notificador.total, 4)
}

func postJSON[T any](t *testing.T, r http.Handler, path string, body any, want int) T {
	t.Helper()
	return doJSON[T](t, r, http.MethodPost, path, body, want)
}

func putJSON[T any](t *testing.T, r http.Handler, path string, body any, want int) T {
	t.Helper()
	return doJSON[T](t, r, http.MethodPut, path, body, want)
}

func getJSON[T any](t *testing.T, r http.Handler, path string, want int) T {
	t.Helper()
	return doJSON[T](t, r, http.MethodGet, path, nil, want)
}

func doJSON[T any](t *testing.T, r http.Handler, method, path string, body any, want int) T {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	require.Equal(t, want, rec.Code, rec.Body.String())
	var out T
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
	return out
}
