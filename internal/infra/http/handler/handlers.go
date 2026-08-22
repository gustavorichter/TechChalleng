package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/techchalleng/oficina/internal/application/dto"
	"github.com/techchalleng/oficina/internal/application/usecase"
)

type AuthHandler struct {
	uc *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// Login godoc
// @Summary      Autenticação de administrador
// @Description  Realiza login e retorna token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Credenciais"
// @Success      200   {object}  dto.LoginResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	resp, err := h.uc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

type ClienteHandler struct {
	uc *usecase.ClienteUseCase
}

func NewClienteHandler(uc *usecase.ClienteUseCase) *ClienteHandler {
	return &ClienteHandler{uc: uc}
}

// Criar godoc
// @Summary      Criar cliente
// @Description  Cadastra um novo cliente (requer JWT admin)
// @Tags         clientes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CriarClienteRequest  true  "Dados do cliente"
// @Success      201   {object}  dto.ClienteResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Router       /clientes [post]
func (h *ClienteHandler) Criar(c *gin.Context) {
	var req dto.CriarClienteRequest
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
// @Summary      Buscar cliente por ID
// @Tags         clientes
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "ID do cliente"
// @Success      200  {object}  dto.ClienteResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /clientes/{id} [get]
func (h *ClienteHandler) Buscar(c *gin.Context) {
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
// @Summary      Listar clientes
// @Tags         clientes
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  dto.ClienteResponse
// @Router       /clientes [get]
func (h *ClienteHandler) Listar(c *gin.Context) {
	resp, err := h.uc.Listar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Excluir godoc
// @Summary      Excluir cliente
// @Tags         clientes
// @Security     BearerAuth
// @Param        id  path  string  true  "ID do cliente"
// @Success      204
// @Router       /clientes/{id} [delete]
func (h *ClienteHandler) Excluir(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	if err := h.uc.Excluir(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

type VeiculoHandler struct {
	uc *usecase.VeiculoUseCase
}

func NewVeiculoHandler(uc *usecase.VeiculoUseCase) *VeiculoHandler {
	return &VeiculoHandler{uc: uc}
}

// Criar godoc
// @Summary      Criar veículo
// @Tags         veiculos
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CriarVeiculoRequest  true  "Dados do veículo"
// @Success      201   {object}  dto.VeiculoResponse
// @Router       /veiculos [post]
func (h *VeiculoHandler) Criar(c *gin.Context) {
	var req dto.CriarVeiculoRequest
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
// @Summary      Buscar veículo por ID
// @Tags         veiculos
// @Security     BearerAuth
// @Param        id  path  string  true  "ID do veículo"
// @Success      200  {object}  dto.VeiculoResponse
// @Router       /veiculos/{id} [get]
func (h *VeiculoHandler) Buscar(c *gin.Context) {
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
// @Summary      Listar veículos
// @Tags         veiculos
// @Security     BearerAuth
// @Success      200  {array}  dto.VeiculoResponse
// @Router       /veiculos [get]
func (h *VeiculoHandler) Listar(c *gin.Context) {
	resp, err := h.uc.Listar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Excluir godoc
// @Summary      Excluir veículo
// @Tags         veiculos
// @Security     BearerAuth
// @Param        id  path  string  true  "ID do veículo"
// @Success      204
// @Router       /veiculos/{id} [delete]
func (h *VeiculoHandler) Excluir(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	if err := h.uc.Excluir(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

type ServicoHandler struct {
	uc *usecase.ServicoUseCase
}

func NewServicoHandler(uc *usecase.ServicoUseCase) *ServicoHandler {
	return &ServicoHandler{uc: uc}
}

// Criar godoc
// @Summary      Criar serviço
// @Tags         servicos
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CriarServicoRequest  true  "Dados do serviço"
// @Success      201   {object}  dto.ServicoResponse
// @Router       /servicos [post]
func (h *ServicoHandler) Criar(c *gin.Context) {
	var req dto.CriarServicoRequest
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
// @Summary      Buscar serviço por ID
// @Tags         servicos
// @Security     BearerAuth
// @Param        id  path  string  true  "ID do serviço"
// @Success      200  {object}  dto.ServicoResponse
// @Router       /servicos/{id} [get]
func (h *ServicoHandler) Buscar(c *gin.Context) {
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
// @Summary      Listar serviços
// @Tags         servicos
// @Security     BearerAuth
// @Success      200  {array}  dto.ServicoResponse
// @Router       /servicos [get]
func (h *ServicoHandler) Listar(c *gin.Context) {
	resp, err := h.uc.Listar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Excluir godoc
// @Summary      Excluir serviço
// @Tags         servicos
// @Security     BearerAuth
// @Param        id  path  string  true  "ID do serviço"
// @Success      204
// @Router       /servicos/{id} [delete]
func (h *ServicoHandler) Excluir(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	if err := h.uc.Excluir(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

type PecaHandler struct {
	uc *usecase.PecaUseCase
}

func NewPecaHandler(uc *usecase.PecaUseCase) *PecaHandler {
	return &PecaHandler{uc: uc}
}

// Criar godoc
// @Summary      Criar peça/insumo
// @Tags         pecas
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CriarPecaRequest  true  "Dados da peça"
// @Success      201   {object}  dto.PecaResponse
// @Router       /pecas [post]
func (h *PecaHandler) Criar(c *gin.Context) {
	var req dto.CriarPecaRequest
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
// @Summary      Buscar peça por ID
// @Tags         pecas
// @Security     BearerAuth
// @Param        id  path  string  true  "ID da peça"
// @Success      200  {object}  dto.PecaResponse
// @Router       /pecas/{id} [get]
func (h *PecaHandler) Buscar(c *gin.Context) {
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
// @Summary      Listar peças
// @Tags         pecas
// @Security     BearerAuth
// @Success      200  {array}  dto.PecaResponse
// @Router       /pecas [get]
func (h *PecaHandler) Listar(c *gin.Context) {
	resp, err := h.uc.Listar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Excluir godoc
// @Summary      Excluir peça
// @Tags         pecas
// @Security     BearerAuth
// @Param        id  path  string  true  "ID da peça"
// @Success      204
// @Router       /pecas/{id} [delete]
func (h *PecaHandler) Excluir(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}
	if err := h.uc.Excluir(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
