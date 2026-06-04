package postgres

import (
	"context"
	"fmt"
	"time"

	entity "final/internal/entities"
	"final/internal/ports"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ ports.PriceRepository = (*PriceRepositoryPostgres)(nil)

type PriceRepositoryPostgres struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

func NewPriceRepositoryPostgres(connString string) (*PriceRepositoryPostgres, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("NewPriceRepositoryPostgres: failed to connect to database: %w", err)
	}

	return &PriceRepositoryPostgres{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

// Close — закрывает пул соединений
func (r *PriceRepositoryPostgres) Close() {
	r.pool.Close()
}

// SavePrices — массовое сохранение цен с проверкой валют
func (r *PriceRepositoryPostgres) SavePrices(ctx context.Context, prices []entity.Price) error {
	if len(prices) == 0 {
		return nil
	}

	// 1. Проверяем и добавляем валюты в таблицу currencies
	for _, p := range prices {
		// Проверяем, существует ли валюта
		var exists bool
		err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM currencies WHERE symbol = $1)", p.Symbol).Scan(&exists)
		if err != nil {
			return fmt.Errorf("SavePrices: failed to check currency %s: %w", p.Symbol, err)
		}

		// Если нет — добавляем
		if !exists {
			_, err := r.pool.Exec(ctx, "INSERT INTO currencies (symbol, name, created_at) VALUES ($1, $2, NOW())", p.Symbol, p.Symbol)
			if err != nil {
				return fmt.Errorf("SavePrices: failed to insert currency %s: %w", p.Symbol, err)
			}
		}
	}

	// 2. Сохраняем цены
	builder := r.sq.Insert("prices").Columns("symbol", "price", "created_at")
	for _, p := range prices {
		builder = builder.Values(p.Symbol, p.Price, p.CreatedAt)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("SavePrices: failed to build insert: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SavePrices: failed to execute insert: %w", err)
	}

	return nil
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
		return nil, fmt.Errorf("GetPricesLast: failed to build SQL query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("GetPricesLast: failed to query prices for symbols %v: %w", symbols, err)
	}
	defer rows.Close()

	var result []entity.Price
	for rows.Next() {
		var p entity.Price
		if err := rows.Scan(&p.Symbol, &p.Price, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetPricesLast: failed to scan row: %w", err)
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetPricesLast: rows iteration error: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("GetPricesLast: no prices found for symbols %v", symbols)
	}

	return result, nil
}

// GetAllSymbols возвращает все символы, которые есть в таблице prices
func (r *PriceRepositoryPostgres) GetAllSymbols(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx, "SELECT DISTINCT symbol FROM prices ORDER BY symbol")
	if err != nil {
		return nil, fmt.Errorf("GetAllSymbols: failed to query symbols: %w", err)
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("GetAllSymbols: failed to scan row: %w", err)
		}
		symbols = append(symbols, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllSymbols: rows iteration error: %w", err)
	}

	return symbols, nil
}

func (r *PriceRepositoryPostgres) GetMinPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	sql := `
        SELECT symbol, MIN(price) as price, MIN(created_at) as created_at
        FROM prices
        WHERE symbol = ANY($1)
        GROUP BY symbol
    `

	rows, err := r.pool.Query(ctx, sql, symbols)
	if err != nil {
		return nil, fmt.Errorf("GetMinPrices: failed to query min prices for symbols %v: %w", symbols, err)
	}
	defer rows.Close()

	var result []entity.Price
	for rows.Next() {
		var p entity.Price
		if err := rows.Scan(&p.Symbol, &p.Price, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetMinPrices: failed to scan row: %w", err)
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetMinPrices: rows iteration error: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("GetMinPrices: no min prices found for symbols %v", symbols)
	}

	return result, nil
}

// GetMaxPrices — максимальные цены для списка валют
func (r *PriceRepositoryPostgres) GetMaxPrices(ctx context.Context, symbols []string) ([]entity.Price, error) {
	sql := `
        SELECT symbol, MAX(price) as price, MAX(created_at) as created_at
        FROM prices
        WHERE symbol = ANY($1)
        GROUP BY symbol
    `

	rows, err := r.pool.Query(ctx, sql, symbols)
	if err != nil {
		return nil, fmt.Errorf("GetMaxPrices: failed to query max prices for symbols %v: %w", symbols, err)
	}
	defer rows.Close()

	var result []entity.Price
	for rows.Next() {
		var p entity.Price
		if err := rows.Scan(&p.Symbol, &p.Price, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetMaxPrices: failed to scan row: %w", err)
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetMaxPrices: rows iteration error: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("GetMaxPrices: no max prices found for symbols %v", symbols)
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

// GetChangePercent возвращает изменение за час
func (r *PriceRepositoryPostgres) GetChangePercent(ctx context.Context, symbols []string) ([]entity.Price, error) {
	var result []entity.Price

	for _, symbol := range symbols {
		var currentPrice, hourAgoPrice float64

		// Текущая цена
		err := r.pool.QueryRow(ctx, `
            SELECT price FROM prices 
            WHERE symbol = $1 
            ORDER BY created_at DESC 
            LIMIT 1
        `, symbol).Scan(&currentPrice)
		if err != nil {
			return nil, fmt.Errorf("GetChangePercent: failed to get current price for %s: %w", symbol, err)
		}

		// Цена час назад
		err = r.pool.QueryRow(ctx, `
            SELECT price FROM prices 
            WHERE symbol = $1 AND created_at <= NOW() - INTERVAL '1 hour'
            ORDER BY created_at DESC 
            LIMIT 1
        `, symbol).Scan(&hourAgoPrice)
		if err != nil {
			return nil, fmt.Errorf("GetChangePercent: failed to get hour ago price for %s: %w", symbol, err)
		}

		if hourAgoPrice == 0 {
			return nil, fmt.Errorf("GetChangePercent: hour ago price for %s is zero, skipping", symbol)
		}

		changePercent := ((currentPrice - hourAgoPrice) / hourAgoPrice) * 100

		result = append(result, entity.Price{
			Symbol: symbol,
			Price:  changePercent,
		})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("GetChangePercent: no data for symbols %v", symbols)
	}

	return result, nil
}

// GetExistingSymbols — возвращает только те символы, которые есть в таблице currencies
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
		return nil, fmt.Errorf("GetExistingSymbols: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("GetExistingSymbols: failed to query: %w", err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("GetExistingSymbols: failed to scan: %w", err)
		}
		result = append(result, s)
	}

	return result, nil
}
