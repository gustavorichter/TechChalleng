package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/techchalleng/oficina/internal/application/dto"
	"github.com/techchalleng/oficina/internal/application/usecase"
)

type OrdemServicoHandler struct {
	uc *usecase.OrdemServicoUseCase
}

func NewOrdemServicoHandler(uc *usecase.OrdemServicoUseCase) *OrdemServicoHandler {
	return &OrdemServicoHandler{uc: uc}
}

// Criar godoc
// @Summary      Criar ordem de serviço
// @Description  Cria uma nova OS no status "Recebida"
// @Tags         ordens-servico
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CriarOSRequest  true  "Dados da OS"
// @Success      201   {object}  dto.OrdemServicoResponse
// @Router       /ordens-servico [post]
func (h *OrdemServicoHandler) Criar(c *gin.Context) {
	var req dto.CriarOSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	resp, err := h.uc.Criar(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Buscar godoc
// @Summary      Buscar ordem de serviço por ID
// @Tags         ordens-servico
// @Security     BearerAuth
// @Param        id  path  string  true  "ID da OS"
// @Success      200  {object}  dto.OrdemServicoResponse
// @Router       /ordens-servico/{id} [get]
func (h *OrdemServicoHandler) Buscar(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	resp, err := h.uc.BuscarPorID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Listar godoc
// @Summary      Listar ordens de serviço
// @Tags         ordens-servico
// @Security     BearerAuth
// @Success      200  {array}  dto.OrdemServicoResponse
// @Router       /ordens-servico [get]
func (h *OrdemServicoHandler) Listar(c *gin.Context) {
	resp, err := h.uc.Listar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// AdicionarServico godoc
// @Summary      Adicionar serviço à OS
// @Tags         ordens-servico
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                      true  "ID da OS"
// @Param        body  body  dto.AdicionarServicoOSRequest  true  "Serviço"
// @Success      200   {object}  dto.OrdemServicoResponse
// @Router       /ordens-servico/{id}/servicos [post]
func (h *OrdemServicoHandler) AdicionarServico(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	var req dto.AdicionarServicoOSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	resp, err := h.uc.AdicionarServico(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// AdicionarPeca godoc
// @Summary      Adicionar peça à OS
// @Description  Verifica saldo de estoque antes de incluir
// @Tags         ordens-servico
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                   true  "ID da OS"
// @Param        body  body  dto.AdicionarPecaOSRequest  true  "Peça"
// @Success      200   {object}  dto.OrdemServicoResponse
// @Router       /ordens-servico/{id}/pecas [post]
func (h *OrdemServicoHandler) AdicionarPeca(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	var req dto.AdicionarPecaOSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	resp, err := h.uc.AdicionarPeca(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// AtualizarStatus godoc
// @Summary      Atualizar status da OS
// @Description  Transição controlada por máquina de estados. Ao "Finalizada", baixa estoque automaticamente.
// @Tags         ordens-servico
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                      true  "ID da OS"
// @Param        body  body  dto.AtualizarStatusOSRequest  true  "Novo status"
// @Success      200   {object}  dto.OrdemServicoResponse
// @Router       /ordens-servico/{id}/status [put]
func (h *OrdemServicoHandler) AtualizarStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	var req dto.AtualizarStatusOSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	resp, err := h.uc.AtualizarStatus(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// AvancarStatus godoc
// @Summary      Avançar status da OS automaticamente
// @Tags         ordens-servico
// @Security     BearerAuth
// @Param        id  path  string  true  "ID da OS"
// @Success      200  {object}  dto.OrdemServicoResponse
// @Router       /ordens-servico/{id}/avancar [post]
func (h *OrdemServicoHandler) AvancarStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	resp, err := h.uc.AvancarStatus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// AcompanharStatus godoc
// @Summary      Consulta pública de status da OS
// @Description  Rota pública para o cliente acompanhar o progresso da OS sem autenticação
// @Tags         ordens-servico
// @Produce      json
// @Param        id  path  string  true  "ID da OS"
// @Success      200  {object}  dto.StatusOSResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /ordens-servico/{id}/status [get]
func (h *OrdemServicoHandler) AcompanharStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	resp, err := h.uc.AcompanharStatus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// TempoMedio godoc
// @Summary      Tempo médio de execução das OS
// @Description  Retorna o tempo médio em horas entre criação e entrega
// @Tags         ordens-servico
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.TempoMedioResponse
// @Router       /ordens-servico/metricas/tempo-medio [get]
func (h *OrdemServicoHandler) TempoMedio(c *gin.Context) {
	resp, err := h.uc.TempoMedioExecucao(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
