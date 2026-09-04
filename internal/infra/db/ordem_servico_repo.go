package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/repository"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

type OrdemServicoRepo struct {
	pool *pgxpool.Pool
}

func NewOrdemServicoRepo(pool *pgxpool.Pool) repository.OrdemServicoRepository {
	return &OrdemServicoRepo{pool: pool}
}

func (r *OrdemServicoRepo) Criar(ctx context.Context, os *entity.OrdemServico) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
		INSERT INTO ordens_servico (id, cliente_id, veiculo_id, status, valor_total, observacoes, criado_em, entregue_em, atualizado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		os.ID, os.ClienteID, os.VeiculoID, string(os.Status), os.ValorTotal, os.Observacoes,
		os.CriadoEm, os.EntregueEm, os.AtualizadoEm)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *OrdemServicoRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.OrdemServico, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, cliente_id, veiculo_id, status, valor_total, observacoes, criado_em, entregue_em, atualizado_em
		FROM ordens_servico WHERE id = $1`, id)

	os, err := scanOrdemServico(row)
	if err != nil {
		return nil, err
	}

	servicos, err := r.loadServicos(ctx, id)
	if err != nil {
		return nil, err
	}
	pecas, err := r.loadPecas(ctx, id)
	if err != nil {
		return nil, err
	}

	os.Servicos = servicos
	os.Pecas = pecas
	return os, nil
}

func (r *OrdemServicoRepo) Listar(ctx context.Context) ([]*entity.OrdemServico, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, cliente_id, veiculo_id, status, valor_total, observacoes, criado_em, entregue_em, atualizado_em
		FROM ordens_servico ORDER BY criado_em DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.OrdemServico
	for rows.Next() {
		os, err := scanOrdemServicoFromRows(rows)
		if err != nil {
			return nil, err
		}
		servicos, _ := r.loadServicos(ctx, os.ID)
		pecas, _ := r.loadPecas(ctx, os.ID)
		os.Servicos = servicos
		os.Pecas = pecas
		result = append(result, os)
	}
	return result, rows.Err()
}

func (r *OrdemServicoRepo) Atualizar(ctx context.Context, os *entity.OrdemServico) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
		UPDATE ordens_servico SET status=$2, valor_total=$3, observacoes=$4, entregue_em=$5, atualizado_em=$6
		WHERE id=$1`,
		os.ID, string(os.Status), os.ValorTotal, os.Observacoes, os.EntregueEm, os.AtualizadoEm)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM os_servicos WHERE ordem_servico_id = $1`, os.ID)
	if err != nil {
		return err
	}
	for _, s := range os.Servicos {
		_, err = tx.Exec(ctx, `
			INSERT INTO os_servicos (ordem_servico_id, servico_id, nome, valor_mao_obra)
			VALUES ($1, $2, $3, $4)`,
			os.ID, s.ServicoID, s.Nome, s.ValorMaoObra)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx, `DELETE FROM os_pecas WHERE ordem_servico_id = $1`, os.ID)
	if err != nil {
		return err
	}
	for _, p := range os.Pecas {
		_, err = tx.Exec(ctx, `
			INSERT INTO os_pecas (ordem_servico_id, peca_id, nome, quantidade, valor_unitario)
			VALUES ($1, $2, $3, $4, $5)`,
			os.ID, p.PecaID, p.Nome, p.Quantidade, p.ValorUnitario)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *OrdemServicoRepo) TempoMedioExecucao(ctx context.Context) (float64, error) {
	var media *float64
	err := r.pool.QueryRow(ctx, `
		SELECT AVG(EXTRACT(EPOCH FROM (entregue_em - criado_em)) / 3600.0)
		FROM ordens_servico WHERE entregue_em IS NOT NULL`).Scan(&media)
	if err != nil {
		return 0, err
	}
	if media == nil {
		return 0, nil
	}
	return *media, nil
}

func (r *OrdemServicoRepo) loadServicos(ctx context.Context, osID uuid.UUID) ([]entity.ItemServicoOS, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT servico_id, nome, valor_mao_obra FROM os_servicos WHERE ordem_servico_id = $1`, osID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.ItemServicoOS
	for rows.Next() {
		var item entity.ItemServicoOS
		if err := rows.Scan(&item.ServicoID, &item.Nome, &item.ValorMaoObra); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *OrdemServicoRepo) loadPecas(ctx context.Context, osID uuid.UUID) ([]entity.ItemPecaOS, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT peca_id, nome, quantidade, valor_unitario FROM os_pecas WHERE ordem_servico_id = $1`, osID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.ItemPecaOS
	for rows.Next() {
		var item entity.ItemPecaOS
		if err := rows.Scan(&item.PecaID, &item.Nome, &item.Quantidade, &item.ValorUnitario); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanOrdemServico(row pgx.Row) (*entity.OrdemServico, error) {
	var os entity.OrdemServico
	var status string
	if err := row.Scan(&os.ID, &os.ClienteID, &os.VeiculoID, &status, &os.ValorTotal,
		&os.Observacoes, &os.CriadoEm, &os.EntregueEm, &os.AtualizadoEm); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("ordem de serviço não encontrada")
		}
		return nil, err
	}
	os.Status = valobj.StatusOS(status)
	os.Servicos = []entity.ItemServicoOS{}
	os.Pecas = []entity.ItemPecaOS{}
	return &os, nil
}

func scanOrdemServicoFromRows(rows pgx.Rows) (*entity.OrdemServico, error) {
	var os entity.OrdemServico
	var status string
	if err := rows.Scan(&os.ID, &os.ClienteID, &os.VeiculoID, &status, &os.ValorTotal,
		&os.Observacoes, &os.CriadoEm, &os.EntregueEm, &os.AtualizadoEm); err != nil {
		return nil, err
	}
	os.Status = valobj.StatusOS(status)
	os.Servicos = []entity.ItemServicoOS{}
	os.Pecas = []entity.ItemPecaOS{}
	return &os, nil
}
