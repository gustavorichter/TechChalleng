package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/techchalleng/oficina/internal/application/usecase"
	"github.com/techchalleng/oficina/internal/infra/http/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handlers agrupa todos os handlers HTTP da aplicação.
type Handlers struct {
	Auth         *AuthHandler
	Cliente      *ClienteHandler
	Veiculo      *VeiculoHandler
	Servico      *ServicoHandler
	Peca         *PecaHandler
	OrdemServico *OrdemServicoHandler
}

func NewHandlers(
	authUC *usecase.AuthUseCase,
	clienteUC *usecase.ClienteUseCase,
	veiculoUC *usecase.VeiculoUseCase,
	servicoUC *usecase.ServicoUseCase,
	pecaUC *usecase.PecaUseCase,
	osUC *usecase.OrdemServicoUseCase,
) *Handlers {
	return &Handlers{
		Auth:         NewAuthHandler(authUC),
		Cliente:      NewClienteHandler(clienteUC),
		Veiculo:      NewVeiculoHandler(veiculoUC),
		Servico:      NewServicoHandler(servicoUC),
		Peca:         NewPecaHandler(pecaUC),
		OrdemServico: NewOrdemServicoHandler(osUC),
	}
}

// RegisterRoutes configura todas as rotas da API.
func RegisterRoutes(r *gin.Engine, h *Handlers, authUC *usecase.AuthUseCase) {
	r.Use(middleware.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// Auth (pública)
	v1.POST("/auth/login", h.Auth.Login)

	// Consulta pública de status da OS
	v1.GET("/ordens-servico/:id/status", h.OrdemServico.AcompanharStatus)

	// Rotas administrativas protegidas por JWT
	admin := v1.Group("")
	admin.Use(middleware.JWTAuth(authUC))
	{
		// Clientes
		admin.POST("/clientes", h.Cliente.Criar)
		admin.GET("/clientes", h.Cliente.Listar)
		admin.GET("/clientes/:id", h.Cliente.Buscar)
		admin.DELETE("/clientes/:id", h.Cliente.Excluir)

		// Veículos
		admin.POST("/veiculos", h.Veiculo.Criar)
		admin.GET("/veiculos", h.Veiculo.Listar)
		admin.GET("/veiculos/:id", h.Veiculo.Buscar)
		admin.DELETE("/veiculos/:id", h.Veiculo.Excluir)

		// Serviços
		admin.POST("/servicos", h.Servico.Criar)
		admin.GET("/servicos", h.Servico.Listar)
		admin.GET("/servicos/:id", h.Servico.Buscar)
		admin.DELETE("/servicos/:id", h.Servico.Excluir)

		// Peças / Estoque
		admin.POST("/pecas", h.Peca.Criar)
		admin.GET("/pecas", h.Peca.Listar)
		admin.GET("/pecas/:id", h.Peca.Buscar)
		admin.DELETE("/pecas/:id", h.Peca.Excluir)

		// Ordens de Serviço
		admin.POST("/ordens-servico", h.OrdemServico.Criar)
		admin.GET("/ordens-servico", h.OrdemServico.Listar)
		admin.GET("/ordens-servico/metricas/tempo-medio", h.OrdemServico.TempoMedio)
		admin.GET("/ordens-servico/:id", h.OrdemServico.Buscar)
		admin.POST("/ordens-servico/:id/servicos", h.OrdemServico.AdicionarServico)
		admin.POST("/ordens-servico/:id/pecas", h.OrdemServico.AdicionarPeca)
		admin.PUT("/ordens-servico/:id/status", h.OrdemServico.AtualizarStatus)
		admin.POST("/ordens-servico/:id/avancar", h.OrdemServico.AvancarStatus)
	}
}
