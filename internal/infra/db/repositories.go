package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/repository"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

type ClienteRepo struct {
	pool *pgxpool.Pool
}

func NewClienteRepo(pool *pgxpool.Pool) repository.ClienteRepository {
	return &ClienteRepo{pool: pool}
}

func (r *ClienteRepo) Criar(ctx context.Context, c *entity.Cliente) error {
	var cpf, cnpj *string
	if c.CPF != nil {
		s := c.CPF.String()
		cpf = &s
	}
	if c.CNPJ != nil {
		s := c.CNPJ.String()
		cnpj = &s
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO clientes (id, nome, cpf, cnpj, email, telefone, criado_em, atualizado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.Nome, cpf, cnpj, c.Email.String(), c.Telefone, c.CriadoEm, c.AtualizadoEm)
	return err
}

func (r *ClienteRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Cliente, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, nome, cpf, cnpj, email, telefone, criado_em, atualizado_em
		FROM clientes WHERE id = $1`, id)
	return scanCliente(row)
}

func (r *ClienteRepo) Listar(ctx context.Context) ([]*entity.Cliente, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, nome, cpf, cnpj, email, telefone, criado_em, atualizado_em
		FROM clientes ORDER BY criado_em DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanClientes(rows)
}

func (r *ClienteRepo) Atualizar(ctx context.Context, c *entity.Cliente) error {
	var cpf, cnpj *string
	if c.CPF != nil {
		s := c.CPF.String()
		cpf = &s
	}
	if c.CNPJ != nil {
		s := c.CNPJ.String()
		cnpj = &s
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE clientes SET nome=$2, cpf=$3, cnpj=$4, email=$5, telefone=$6, atualizado_em=$7
		WHERE id=$1`, c.ID, c.Nome, cpf, cnpj, c.Email.String(), c.Telefone, c.AtualizadoEm)
	return err
}

func (r *ClienteRepo) Excluir(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM clientes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("cliente não encontrado")
	}
	return nil
}

func scanCliente(row pgx.Row) (*entity.Cliente, error) {
	var c entity.Cliente
	var cpf, cnpj *string
	var email string
	if err := row.Scan(&c.ID, &c.Nome, &cpf, &cnpj, &email, &c.Telefone, &c.CriadoEm, &c.AtualizadoEm); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("cliente não encontrado")
		}
		return nil, err
	}
	em, err := valobj.NewEmail(email)
	if err != nil {
		return nil, err
	}
	c.Email = em
	if cpf != nil {
		v, err := valobj.NewCPF(*cpf)
		if err != nil {
			return nil, err
		}
		c.CPF = &v
	}
	if cnpj != nil {
		v, err := valobj.NewCNPJ(*cnpj)
		if err != nil {
			return nil, err
		}
		c.CNPJ = &v
	}
	return &c, nil
}

func scanClientes(rows pgx.Rows) ([]*entity.Cliente, error) {
	var result []*entity.Cliente
	for rows.Next() {
		var c entity.Cliente
		var cpf, cnpj *string
		var email string
		if err := rows.Scan(&c.ID, &c.Nome, &cpf, &cnpj, &email, &c.Telefone, &c.CriadoEm, &c.AtualizadoEm); err != nil {
			return nil, err
		}
		em, err := valobj.NewEmail(email)
		if err != nil {
			return nil, err
		}
		c.Email = em
		if cpf != nil {
			v, _ := valobj.NewCPF(*cpf)
			c.CPF = &v
		}
		if cnpj != nil {
			v, _ := valobj.NewCNPJ(*cnpj)
			c.CNPJ = &v
		}
		result = append(result, &c)
	}
	return result, rows.Err()
}

type VeiculoRepo struct {
	pool *pgxpool.Pool
}

func NewVeiculoRepo(pool *pgxpool.Pool) repository.VeiculoRepository {
	return &VeiculoRepo{pool: pool}
}

func (r *VeiculoRepo) Criar(ctx context.Context, v *entity.Veiculo) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO veiculos (id, cliente_id, placa, marca, modelo, ano, criado_em, atualizado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		v.ID, v.ClienteID, v.Placa.String(), v.Marca, v.Modelo, v.Ano, v.CriadoEm, v.AtualizadoEm)
	return err
}

func (r *VeiculoRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Veiculo, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, cliente_id, placa, marca, modelo, ano, criado_em, atualizado_em
		FROM veiculos WHERE id = $1`, id)
	return scanVeiculo(row)
}

func (r *VeiculoRepo) ListarPorCliente(ctx context.Context, clienteID uuid.UUID) ([]*entity.Veiculo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, cliente_id, placa, marca, modelo, ano, criado_em, atualizado_em
		FROM veiculos WHERE cliente_id = $1 ORDER BY criado_em DESC`, clienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVeiculos(rows)
}

func (r *VeiculoRepo) Listar(ctx context.Context) ([]*entity.Veiculo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, cliente_id, placa, marca, modelo, ano, criado_em, atualizado_em
		FROM veiculos ORDER BY criado_em DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVeiculos(rows)
}

func (r *VeiculoRepo) Atualizar(ctx context.Context, v *entity.Veiculo) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE veiculos SET placa=$2, marca=$3, modelo=$4, ano=$5, atualizado_em=$6
		WHERE id=$1`, v.ID, v.Placa.String(), v.Marca, v.Modelo, v.Ano, v.AtualizadoEm)
	return err
}

func (r *VeiculoRepo) Excluir(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM veiculos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("veículo não encontrado")
	}
	return nil
}

func scanVeiculo(row pgx.Row) (*entity.Veiculo, error) {
	var v entity.Veiculo
	var placa string
	if err := row.Scan(&v.ID, &v.ClienteID, &placa, &v.Marca, &v.Modelo, &v.Ano, &v.CriadoEm, &v.AtualizadoEm); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("veículo não encontrado")
		}
		return nil, err
	}
	p, err := valobj.NewPlaca(placa)
	if err != nil {
		return nil, err
	}
	v.Placa = p
	return &v, nil
}

func scanVeiculos(rows pgx.Rows) ([]*entity.Veiculo, error) {
	var result []*entity.Veiculo
	for rows.Next() {
		var v entity.Veiculo
		var placa string
		if err := rows.Scan(&v.ID, &v.ClienteID, &placa, &v.Marca, &v.Modelo, &v.Ano, &v.CriadoEm, &v.AtualizadoEm); err != nil {
			return nil, err
		}
		p, _ := valobj.NewPlaca(placa)
		v.Placa = p
		result = append(result, &v)
	}
	return result, rows.Err()
}

type ServicoRepo struct {
	pool *pgxpool.Pool
}

func NewServicoRepo(pool *pgxpool.Pool) repository.ServicoRepository {
	return &ServicoRepo{pool: pool}
}

func (r *ServicoRepo) Criar(ctx context.Context, s *entity.Servico) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO servicos (id, nome, descricao, valor_mao_obra, ativo, criado_em, atualizado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		s.ID, s.Nome, s.Descricao, s.ValorMaoObra, s.Ativo, s.CriadoEm, s.AtualizadoEm)
	return err
}

func (r *ServicoRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Servico, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, nome, descricao, valor_mao_obra, ativo, criado_em, atualizado_em
		FROM servicos WHERE id = $1`, id)
	var s entity.Servico
	if err := row.Scan(&s.ID, &s.Nome, &s.Descricao, &s.ValorMaoObra, &s.Ativo, &s.CriadoEm, &s.AtualizadoEm); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("serviço não encontrado")
		}
		return nil, err
	}
	return &s, nil
}

func (r *ServicoRepo) Listar(ctx context.Context) ([]*entity.Servico, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, nome, descricao, valor_mao_obra, ativo, criado_em, atualizado_em
		FROM servicos ORDER BY nome`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*entity.Servico
	for rows.Next() {
		var s entity.Servico
		if err := rows.Scan(&s.ID, &s.Nome, &s.Descricao, &s.ValorMaoObra, &s.Ativo, &s.CriadoEm, &s.AtualizadoEm); err != nil {
			return nil, err
		}
		result = append(result, &s)
	}
	return result, rows.Err()
}

func (r *ServicoRepo) Atualizar(ctx context.Context, s *entity.Servico) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE servicos SET nome=$2, descricao=$3, valor_mao_obra=$4, ativo=$5, atualizado_em=$6
		WHERE id=$1`, s.ID, s.Nome, s.Descricao, s.ValorMaoObra, s.Ativo, s.AtualizadoEm)
	return err
}

func (r *ServicoRepo) Excluir(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM servicos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("serviço não encontrado")
	}
	return nil
}

type PecaRepo struct {
	pool *pgxpool.Pool
}

func NewPecaRepo(pool *pgxpool.Pool) repository.PecaRepository {
	return &PecaRepo{pool: pool}
}

func (r *PecaRepo) Criar(ctx context.Context, p *entity.Peca) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO pecas (id, nome, codigo, valor_unitario, quantidade_estoque, ativo, criado_em, atualizado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		p.ID, p.Nome, p.Codigo, p.ValorUnitario, p.QuantidadeEstoque, p.Ativo, p.CriadoEm, p.AtualizadoEm)
	return err
}

func (r *PecaRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Peca, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, nome, codigo, valor_unitario, quantidade_estoque, ativo, criado_em, atualizado_em
		FROM pecas WHERE id = $1`, id)
	var p entity.Peca
	if err := row.Scan(&p.ID, &p.Nome, &p.Codigo, &p.ValorUnitario, &p.QuantidadeEstoque, &p.Ativo, &p.CriadoEm, &p.AtualizadoEm); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("peça %s não encontrada", id)
		}
		return nil, err
	}
	return &p, nil
}

func (r *PecaRepo) Listar(ctx context.Context) ([]*entity.Peca, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, nome, codigo, valor_unitario, quantidade_estoque, ativo, criado_em, atualizado_em
		FROM pecas ORDER BY nome`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*entity.Peca
	for rows.Next() {
		var p entity.Peca
		if err := rows.Scan(&p.ID, &p.Nome, &p.Codigo, &p.ValorUnitario, &p.QuantidadeEstoque, &p.Ativo, &p.CriadoEm, &p.AtualizadoEm); err != nil {
			return nil, err
		}
		result = append(result, &p)
	}
	return result, rows.Err()
}

func (r *PecaRepo) Atualizar(ctx context.Context, p *entity.Peca) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE pecas SET nome=$2, codigo=$3, valor_unitario=$4, quantidade_estoque=$5, ativo=$6, atualizado_em=$7
		WHERE id=$1`, p.ID, p.Nome, p.Codigo, p.ValorUnitario, p.QuantidadeEstoque, p.Ativo, p.AtualizadoEm)
	return err
}

func (r *PecaRepo) Excluir(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM pecas WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("peça não encontrada")
	}
	return nil
}
