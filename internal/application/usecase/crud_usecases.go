package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/techchalleng/oficina/internal/application/dto"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/repository"
	"github.com/techchalleng/oficina/internal/domain/valobj"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	jwtSecret []byte
	jwtExpiry time.Duration
}

func NewAuthUseCase() *AuthUseCase {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "oficina-secret-key-change-in-production"
	}
	expiry := 24 * time.Hour
	return &AuthUseCase{
		jwtSecret: []byte(secret),
		jwtExpiry: expiry,
	}
}

type jwtClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (uc *AuthUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	adminUser := os.Getenv("ADMIN_USERNAME")
	if adminUser == "" {
		adminUser = "admin"
	}
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminPass == "" {
		adminPass = "admin123"
	}

	if req.Username != adminUser || req.Password != adminPass {
		return nil, errors.New("credenciais inválidas")
	}

	now := time.Now()
	claims := jwtClaims{
		Username: adminUser,
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(uc.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   adminUser,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar token: %w", err)
	}

	return &dto.LoginResponse{
		Token:     tokenStr,
		ExpiresIn: int(uc.jwtExpiry.Seconds()),
	}, nil
}

func (uc *AuthUseCase) ValidateToken(tokenStr string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado")
		}
		return uc.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}
	return claims, nil
}

// HashPassword utilitário para gerar hash de senha (uso em setup).
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// --- Cliente UseCase ---

type ClienteUseCase struct {
	repo repository.ClienteRepository
}

func NewClienteUseCase(repo repository.ClienteRepository) *ClienteUseCase {
	return &ClienteUseCase{repo: repo}
}

func (uc *ClienteUseCase) Criar(ctx context.Context, req dto.CriarClienteRequest) (*dto.ClienteResponse, error) {
	email, err := valobj.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	cliente := entity.NewCliente(req.Nome, email, req.Telefone)

	if req.CPF != "" {
		cpf, err := valobj.NewCPF(req.CPF)
		if err != nil {
			return nil, err
		}
		cliente.DefinirCPF(cpf)
	} else if req.CNPJ != "" {
		cnpj, err := valobj.NewCNPJ(req.CNPJ)
		if err != nil {
			return nil, err
		}
		cliente.DefinirCNPJ(cnpj)
	} else {
		return nil, errors.New("CPF ou CNPJ é obrigatório")
	}

	if err := uc.repo.Criar(ctx, cliente); err != nil {
		return nil, err
	}
	resp := dto.ToClienteResponse(cliente)
	return &resp, nil
}

func (uc *ClienteUseCase) BuscarPorID(ctx context.Context, id uuid.UUID) (*dto.ClienteResponse, error) {
	cliente, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToClienteResponse(cliente)
	return &resp, nil
}

func (uc *ClienteUseCase) Listar(ctx context.Context) ([]dto.ClienteResponse, error) {
	clientes, err := uc.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ClienteResponse, len(clientes))
	for i, c := range clientes {
		result[i] = dto.ToClienteResponse(c)
	}
	return result, nil
}

func (uc *ClienteUseCase) Excluir(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Excluir(ctx, id)
}

// --- Veículo UseCase ---

type VeiculoUseCase struct {
	veiculoRepo repository.VeiculoRepository
	clienteRepo repository.ClienteRepository
}

func NewVeiculoUseCase(vRepo repository.VeiculoRepository, cRepo repository.ClienteRepository) *VeiculoUseCase {
	return &VeiculoUseCase{veiculoRepo: vRepo, clienteRepo: cRepo}
}

func (uc *VeiculoUseCase) Criar(ctx context.Context, req dto.CriarVeiculoRequest) (*dto.VeiculoResponse, error) {
	if _, err := uc.clienteRepo.BuscarPorID(ctx, req.ClienteID); err != nil {
		return nil, fmt.Errorf("cliente não encontrado: %w", err)
	}

	placa, err := valobj.NewPlaca(req.Placa)
	if err != nil {
		return nil, err
	}

	veiculo := entity.NewVeiculo(req.ClienteID, placa, req.Marca, req.Modelo, req.Ano)
	if err := uc.veiculoRepo.Criar(ctx, veiculo); err != nil {
		return nil, err
	}
	resp := dto.ToVeiculoResponse(veiculo)
	return &resp, nil
}

func (uc *VeiculoUseCase) BuscarPorID(ctx context.Context, id uuid.UUID) (*dto.VeiculoResponse, error) {
	veiculo, err := uc.veiculoRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToVeiculoResponse(veiculo)
	return &resp, nil
}

func (uc *VeiculoUseCase) Listar(ctx context.Context) ([]dto.VeiculoResponse, error) {
	veiculos, err := uc.veiculoRepo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.VeiculoResponse, len(veiculos))
	for i, v := range veiculos {
		result[i] = dto.ToVeiculoResponse(v)
	}
	return result, nil
}

func (uc *VeiculoUseCase) Excluir(ctx context.Context, id uuid.UUID) error {
	return uc.veiculoRepo.Excluir(ctx, id)
}

// --- Serviço UseCase ---

type ServicoUseCase struct {
	repo repository.ServicoRepository
}

func NewServicoUseCase(repo repository.ServicoRepository) *ServicoUseCase {
	return &ServicoUseCase{repo: repo}
}

func (uc *ServicoUseCase) Criar(ctx context.Context, req dto.CriarServicoRequest) (*dto.ServicoResponse, error) {
	if req.ValorMaoObra.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("valor da mão de obra deve ser positivo")
	}
	servico := entity.NewServico(req.Nome, req.Descricao, req.ValorMaoObra)
	if err := uc.repo.Criar(ctx, servico); err != nil {
		return nil, err
	}
	resp := dto.ToServicoResponse(servico)
	return &resp, nil
}

func (uc *ServicoUseCase) BuscarPorID(ctx context.Context, id uuid.UUID) (*dto.ServicoResponse, error) {
	servico, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToServicoResponse(servico)
	return &resp, nil
}

func (uc *ServicoUseCase) Listar(ctx context.Context) ([]dto.ServicoResponse, error) {
	servicos, err := uc.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ServicoResponse, len(servicos))
	for i, s := range servicos {
		result[i] = dto.ToServicoResponse(s)
	}
	return result, nil
}

func (uc *ServicoUseCase) Excluir(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Excluir(ctx, id)
}

// --- Peça UseCase ---

type PecaUseCase struct {
	repo repository.PecaRepository
}

func NewPecaUseCase(repo repository.PecaRepository) *PecaUseCase {
	return &PecaUseCase{repo: repo}
}

func (uc *PecaUseCase) Criar(ctx context.Context, req dto.CriarPecaRequest) (*dto.PecaResponse, error) {
	if req.ValorUnitario.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("valor unitário deve ser positivo")
	}
	peca := entity.NewPeca(req.Nome, req.Codigo, req.ValorUnitario, req.QuantidadeEstoque)
	if err := uc.repo.Criar(ctx, peca); err != nil {
		return nil, err
	}
	resp := dto.ToPecaResponse(peca)
	return &resp, nil
}

func (uc *PecaUseCase) BuscarPorID(ctx context.Context, id uuid.UUID) (*dto.PecaResponse, error) {
	peca, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToPecaResponse(peca)
	return &resp, nil
}

func (uc *PecaUseCase) Listar(ctx context.Context) ([]dto.PecaResponse, error) {
	pecas, err := uc.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.PecaResponse, len(pecas))
	for i, p := range pecas {
		result[i] = dto.ToPecaResponse(p)
	}
	return result, nil
}

func (uc *PecaUseCase) Excluir(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Excluir(ctx, id)
}
