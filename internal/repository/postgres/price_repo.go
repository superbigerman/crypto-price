package postgres

import (
	"context"
	"fmt"
	"time"

	"final/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PriceRepositoryPostgres struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

func NewPriceRepositoryPostgres(pool *pgxpool.Pool) *PriceRepositoryPostgres {
	return &PriceRepositoryPostgres{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// BuildConnString собирает строку подключения из конфига
func BuildConnString(host, port, user, password, dbname, sslmode string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}

// SavePrices — массовое сохранение цен
func (r *PriceRepositoryPostgres) SavePrices(ctx context.Context, prices []entity.Price) error {
	if len(prices) == 0 {
		return nil
	}

	builder := r.sq.Insert("prices").Columns("symbol", "price", "created_at")
	for _, p := range prices {
		builder = builder.Values(p.Symbol, p.Price, p.CreatedAt)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	return err
}

// GetPricesLast — последние цены для списка валют
func (r *PriceRepositoryPostgres) GetPricesLast(ctx context.Context, symbols []string) ([]entity.Price, error) {
	sql, args, err := r.sq.
		Select("DISTINCT ON (symbol) symbol", "price", "created_at").
		From("prices").
		Where(squirrel.Eq{"symbol": symbols}).
		OrderBy("symbol", "created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Price
	for rows.Next() {
		var p entity.Price
		if err := rows.Scan(&p.Symbol, &p.Price, &p.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

// GetMinPrices — минимальные цены для списка валют
func (r *PriceRepositoryPostgres) GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	sql, args, err := r.sq.
		Select("DISTINCT ON (symbol) symbol", "price", "created_at").
		From("prices").
		Where(squirrel.Eq{"symbol": symbols}).
		OrderBy("symbol", "price ASC", "created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Price
	for rows.Next() {
		var p entity.Price
		if err := rows.Scan(&p.Symbol, &p.Price, &p.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

// GetMaxPrices — максимальные цены для списка валют
func (r *PriceRepositoryPostgres) GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	sql, args, err := r.sq.
		Select("DISTINCT ON (symbol) symbol", "price", "created_at").
		From("prices").
		Where(squirrel.Eq{"symbol": symbols}).
		OrderBy("symbol", "price DESC", "created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Price
	for rows.Next() {
		var p entity.Price
		if err := rows.Scan(&p.Symbol, &p.Price, &p.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

// GetExistingSymbols — возвращает символы, которые есть в таблице currencies
func (r *PriceRepositoryPostgres) GetExistingSymbols(ctx context.Context, symbols []string) ([]string, error) {
	if len(symbols) == 0 {
		return []string{}, nil
	}

	sql, args, err := r.sq.
		Select("DISTINCT symbol").
		From("currencies").
		Where(squirrel.Eq{"symbol": symbols}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

// AddCurrency — добавляет новую валюту в таблицу currencies
func (r *PriceRepositoryPostgres) AddCurrency(ctx context.Context, symbol string) error {
	sql, args, err := r.sq.
		Insert("currencies").
		Columns("symbol", "name", "created_at").
		Values(symbol, symbol, time.Now()).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		// Игнорируем ошибку уникальности (валюта уже есть)
		return nil
	}
	return nil
}

// GetChangePercent — заглушка, будет реализовано позже
func (r *PriceRepositoryPostgres) GetChangePercent(ctx context.Context, symbols []string) ([]float64, error) {
	return nil, fmt.Errorf("not implemented yet")
}
