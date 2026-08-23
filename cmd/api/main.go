// Package main Sistema Integrado de Atendimento e Execução de Serviços - Oficina Mecânica
//
// @title           Oficina Mecânica API
// @version         1.0
// @description     API REST para gestão de ordens de serviço, clientes, veículos, serviços e estoque de peças.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   TechChalleng
// @contact.email  suporte@oficina.com
//
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Insira o token JWT no formato: Bearer {token}
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/techchalleng/oficina/internal/application/usecase"
	"github.com/techchalleng/oficina/internal/infra/db"
	"github.com/techchalleng/oficina/internal/infra/http/handler"

	_ "github.com/techchalleng/oficina/docs"
)

func main() {
	ctx := context.Background()

	cfg := db.ConfigFromEnv()
	pool, err := db.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatalf("Erro ao executar migrations: %v", err)
	}
	log.Println("Migrations executadas com sucesso")

	// Repositórios
	clienteRepo := db.NewClienteRepo(pool)
	veiculoRepo := db.NewVeiculoRepo(pool)
	servicoRepo := db.NewServicoRepo(pool)
	pecaRepo := db.NewPecaRepo(pool)
	osRepo := db.NewOrdemServicoRepo(pool)

	// Casos de uso
	authUC := usecase.NewAuthUseCase()
	clienteUC := usecase.NewClienteUseCase(clienteRepo)
	veiculoUC := usecase.NewVeiculoUseCase(veiculoRepo, clienteRepo)
	servicoUC := usecase.NewServicoUseCase(servicoRepo)
	pecaUC := usecase.NewPecaUseCase(pecaRepo)
	osUC := usecase.NewOrdemServicoUseCase(osRepo, clienteRepo, veiculoRepo, servicoRepo, pecaRepo)

	// HTTP
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	handlers := handler.NewHandlers(authUC, clienteUC, veiculoUC, servicoUC, pecaUC, osUC)
	handler.RegisterRoutes(r, handlers, authUC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("Servidor iniciado na porta %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Encerrando servidor...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Erro ao encerrar servidor: %v", err)
	}
}
