package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techchalleng/oficina/internal/application/dto"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

type mockOSRepo struct {
	ordens map[uuid.UUID]*entity.OrdemServico
}

func newMockOSRepo() *mockOSRepo {
	return &mockOSRepo{ordens: make(map[uuid.UUID]*entity.OrdemServico)}
}

func (m *mockOSRepo) Criar(ctx context.Context, os *entity.OrdemServico) error {
	m.ordens[os.ID] = os
	return nil
}

func (m *mockOSRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.OrdemServico, error) {
	os, ok := m.ordens[id]
	if !ok {
		return nil, assert.AnError
	}
	return os, nil
}

func (m *mockOSRepo) Listar(ctx context.Context) ([]*entity.OrdemServico, error) { return nil, nil }
func (m *mockOSRepo) Atualizar(ctx context.Context, os *entity.OrdemServico) error {
	m.ordens[os.ID] = os
	return nil
}
func (m *mockOSRepo) TempoMedioExecucao(ctx context.Context) (float64, error) { return 0, nil }

type mockClienteRepo struct {
	clientes map[uuid.UUID]*entity.Cliente
}

func newMockClienteRepo() *mockClienteRepo {
	return &mockClienteRepo{clientes: make(map[uuid.UUID]*entity.Cliente)}
}

func (m *mockClienteRepo) Criar(ctx context.Context, c *entity.Cliente) error {
	m.clientes[c.ID] = c
	return nil
}
func (m *mockClienteRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Cliente, error) {
	c, ok := m.clientes[id]
	if !ok {
		return nil, assert.AnError
	}
	return c, nil
}
func (m *mockClienteRepo) Listar(ctx context.Context) ([]*entity.Cliente, error) {
	result := make([]*entity.Cliente, 0, len(m.clientes))
	for _, c := range m.clientes {
		result = append(result, c)
	}
	return result, nil
}
func (m *mockClienteRepo) Atualizar(ctx context.Context, c *entity.Cliente) error { return nil }
func (m *mockClienteRepo) Excluir(ctx context.Context, id uuid.UUID) error {
	delete(m.clientes, id)
	return nil
}

type mockVeiculoRepo struct {
	veiculos map[uuid.UUID]*entity.Veiculo
}

func newMockVeiculoRepo() *mockVeiculoRepo {
	return &mockVeiculoRepo{veiculos: make(map[uuid.UUID]*entity.Veiculo)}
}

func (m *mockVeiculoRepo) Criar(ctx context.Context, v *entity.Veiculo) error {
	m.veiculos[v.ID] = v
	return nil
}
func (m *mockVeiculoRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Veiculo, error) {
	v, ok := m.veiculos[id]
	if !ok {
		return nil, assert.AnError
	}
	return v, nil
}
func (m *mockVeiculoRepo) ListarPorCliente(ctx context.Context, id uuid.UUID) ([]*entity.Veiculo, error) {
	return nil, nil
}
func (m *mockVeiculoRepo) Listar(ctx context.Context) ([]*entity.Veiculo, error)  { return nil, nil }
func (m *mockVeiculoRepo) Atualizar(ctx context.Context, v *entity.Veiculo) error { return nil }
func (m *mockVeiculoRepo) Excluir(ctx context.Context, id uuid.UUID) error        { return nil }

type mockServicoRepo struct {
	servicos map[uuid.UUID]*entity.Servico
}

func newMockServicoRepo() *mockServicoRepo {
	return &mockServicoRepo{servicos: make(map[uuid.UUID]*entity.Servico)}
}

func (m *mockServicoRepo) Criar(ctx context.Context, s *entity.Servico) error {
	m.servicos[s.ID] = s
	return nil
}
func (m *mockServicoRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Servico, error) {
	s, ok := m.servicos[id]
	if !ok {
		return nil, assert.AnError
	}
	return s, nil
}
func (m *mockServicoRepo) Listar(ctx context.Context) ([]*entity.Servico, error)  { return nil, nil }
func (m *mockServicoRepo) Atualizar(ctx context.Context, s *entity.Servico) error { return nil }
func (m *mockServicoRepo) Excluir(ctx context.Context, id uuid.UUID) error        { return nil }

type mockPecaRepoUC struct {
	pecas map[uuid.UUID]*entity.Peca
}

func newMockPecaRepoUC() *mockPecaRepoUC {
	return &mockPecaRepoUC{pecas: make(map[uuid.UUID]*entity.Peca)}
}

func (m *mockPecaRepoUC) Criar(ctx context.Context, p *entity.Peca) error {
	m.pecas[p.ID] = p
	return nil
}
func (m *mockPecaRepoUC) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Peca, error) {
	p, ok := m.pecas[id]
	if !ok {
		return nil, assert.AnError
	}
	return p, nil
}
func (m *mockPecaRepoUC) Listar(ctx context.Context) ([]*entity.Peca, error) { return nil, nil }
func (m *mockPecaRepoUC) Atualizar(ctx context.Context, p *entity.Peca) error {
	m.pecas[p.ID] = p
	return nil
}
func (m *mockPecaRepoUC) Excluir(ctx context.Context, id uuid.UUID) error { return nil }

func setupOrdemServicoUC() (*OrdemServicoUseCase, uuid.UUID, uuid.UUID, uuid.UUID) {
	clienteID := uuid.New()
	veiculoID := uuid.New()

	clienteRepo := newMockClienteRepo()
	clienteRepo.clientes[clienteID] = &entity.Cliente{ID: clienteID}

	veiculoRepo := newMockVeiculoRepo()
	veiculoRepo.veiculos[veiculoID] = &entity.Veiculo{ID: veiculoID, ClienteID: clienteID}

	servico := entity.NewServico("Troca de óleo", "Completa", decimal.NewFromFloat(150))
	servicoRepo := newMockServicoRepo()
	servicoRepo.servicos[servico.ID] = servico

	peca := entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(45), 10)
	pecaRepo := newMockPecaRepoUC()
	pecaRepo.pecas[peca.ID] = peca

	osRepo := newMockOSRepo()
	uc := NewOrdemServicoUseCase(osRepo, clienteRepo, veiculoRepo, servicoRepo, pecaRepo)

	return uc, clienteID, veiculoID, servico.ID
}

func TestOrdemServicoUseCase_Criar(t *testing.T) {
	uc, clienteID, veiculoID, _ := setupOrdemServicoUC()

	resp, err := uc.Criar(context.Background(), dto.CriarOSRequest{
		ClienteID: clienteID,
		VeiculoID: veiculoID,
	})
	require.NoError(t, err)
	assert.Equal(t, string(valobj.StatusRecebida), resp.Status)
}

func TestOrdemServicoUseCase_AtualizarStatus_Invalido(t *testing.T) {
	uc, clienteID, veiculoID, _ := setupOrdemServicoUC()

	created, err := uc.Criar(context.Background(), dto.CriarOSRequest{
		ClienteID: clienteID,
		VeiculoID: veiculoID,
	})
	require.NoError(t, err)

	_, err = uc.AtualizarStatus(context.Background(), created.ID, dto.AtualizarStatusOSRequest{
		Status: string(valobj.StatusFinalizada),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "transição inválida")
}

func TestOrdemServicoUseCase_AdicionarServico_RecalculaOrcamento(t *testing.T) {
	uc, clienteID, veiculoID, servicoID := setupOrdemServicoUC()

	created, err := uc.Criar(context.Background(), dto.CriarOSRequest{
		ClienteID: clienteID,
		VeiculoID: veiculoID,
	})
	require.NoError(t, err)

	resp, err := uc.AdicionarServico(context.Background(), created.ID, dto.AdicionarServicoOSRequest{
		ServicoID: servicoID,
	})
	require.NoError(t, err)
	assert.True(t, resp.ValorTotal.Equal(decimal.NewFromFloat(150)))
}

func TestOrdemServicoUseCase_Finalizar_BaixaEstoque(t *testing.T) {
	_, clienteID, veiculoID, servicoID := setupOrdemServicoUC()
	pecaRepo := newMockPecaRepoUC()
	peca := entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(45), 10)
	pecaRepo.pecas[peca.ID] = peca

	clienteRepo := newMockClienteRepo()
	clienteRepo.clientes[clienteID] = &entity.Cliente{ID: clienteID}
	veiculoRepo := newMockVeiculoRepo()
	veiculoRepo.veiculos[veiculoID] = &entity.Veiculo{ID: veiculoID, ClienteID: clienteID}
	servicoRepo := newMockServicoRepo()
	servicoRepo.servicos[servicoID] = entity.NewServico("Troca", "desc", decimal.NewFromFloat(100))
	osRepo := newMockOSRepo()

	uc := NewOrdemServicoUseCase(osRepo, clienteRepo, veiculoRepo, servicoRepo, pecaRepo)

	created, err := uc.Criar(context.Background(), dto.CriarOSRequest{ClienteID: clienteID, VeiculoID: veiculoID})
	require.NoError(t, err)
	_, err = uc.AdicionarPeca(context.Background(), created.ID, dto.AdicionarPecaOSRequest{PecaID: peca.ID, Quantidade: 2})
	require.NoError(t, err)

	statuses := []valobj.StatusOS{
		valobj.StatusEmDiagnostico,
		valobj.StatusAguardandoAprovacao,
		valobj.StatusEmExecucao,
		valobj.StatusFinalizada,
	}
	for _, s := range statuses {
		_, err := uc.AtualizarStatus(context.Background(), created.ID, dto.AtualizarStatusOSRequest{Status: string(s)})
		require.NoError(t, err, "status %s", s)
	}

	assert.Equal(t, 8, pecaRepo.pecas[peca.ID].QuantidadeEstoque)
}

func TestAuthUseCase_Login(t *testing.T) {
	uc := NewAuthUseCase()
	resp, err := uc.Login(context.Background(), dto.LoginRequest{Username: "admin", Password: "admin123"})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)

	claims, err := uc.ValidateToken(resp.Token)
	require.NoError(t, err)
	assert.Equal(t, "admin", claims.Role)
}

func TestAuthUseCase_Login_Invalido(t *testing.T) {
	uc := NewAuthUseCase()
	_, err := uc.Login(context.Background(), dto.LoginRequest{Username: "admin", Password: "wrong"})
	require.Error(t, err)
}
