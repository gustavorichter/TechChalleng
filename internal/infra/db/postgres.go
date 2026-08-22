package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config contém parâmetros de conexão com o PostgreSQL.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func ConfigFromEnv() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "oficina"),
		Password: getEnv("DB_PASSWORD", "oficina_secret"),
		DBName:   getEnv("DB_NAME", "oficina_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// NewPool cria um pool de conexões aguardando o banco estar disponível.
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	dsn := cfg.DSN()
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 30; i++ {
		pool, err = pgxpool.New(ctx, dsn)
		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				return pool, nil
			}
			pool.Close()
		}
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("falha ao conectar ao banco após retries: %w", err)
}

// Migrate executa o schema inicial do banco de dados.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	schema := `
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS clientes (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(11),
    cnpj VARCHAR(14),
    email VARCHAR(255) NOT NULL,
    telefone VARCHAR(20) NOT NULL,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS veiculos (
    id UUID PRIMARY KEY,
    cliente_id UUID NOT NULL REFERENCES clientes(id) ON DELETE CASCADE,
    placa VARCHAR(8) NOT NULL UNIQUE,
    marca VARCHAR(100) NOT NULL,
    modelo VARCHAR(100) NOT NULL,
    ano INT NOT NULL,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS servicos (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    valor_mao_obra NUMERIC(12,2) NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pecas (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    codigo VARCHAR(50) NOT NULL UNIQUE,
    valor_unitario NUMERIC(12,2) NOT NULL,
    quantidade_estoque INT NOT NULL DEFAULT 0,
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ordens_servico (
    id UUID PRIMARY KEY,
    cliente_id UUID NOT NULL REFERENCES clientes(id),
    veiculo_id UUID NOT NULL REFERENCES veiculos(id),
    status VARCHAR(50) NOT NULL DEFAULT 'Recebida',
    valor_total NUMERIC(12,2) NOT NULL DEFAULT 0,
    observacoes TEXT,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    entregue_em TIMESTAMPTZ,
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS os_servicos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ordem_servico_id UUID NOT NULL REFERENCES ordens_servico(id) ON DELETE CASCADE,
    servico_id UUID NOT NULL REFERENCES servicos(id),
    nome VARCHAR(255) NOT NULL,
    valor_mao_obra NUMERIC(12,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS os_pecas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ordem_servico_id UUID NOT NULL REFERENCES ordens_servico(id) ON DELETE CASCADE,
    peca_id UUID NOT NULL REFERENCES pecas(id),
    nome VARCHAR(255) NOT NULL,
    quantidade INT NOT NULL,
    valor_unitario NUMERIC(12,2) NOT NULL
);
`
	_, err := pool.Exec(ctx, schema)
	return err
}
