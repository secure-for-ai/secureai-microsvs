package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"

	//"github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
)

//user=jack password=secret host=pg.example.com port=5432 dbname=mydb
type PGPoolConf struct {
	Host   string `json:"Host"`
	Port   string `json:"Port"`
	DBName string `json:"DBName"`
	User   string `json:"User"`
	PW     string `json:"PW"`
}

type PGError string

func (e PGError) Error() string { return string(e) }

func (c PGPoolConf) GetPGDSN() string {
	pgDSN := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		c.User, c.PW, c.Host, c.Port, c.DBName)
	return pgDSN
}

func NewPGClient(conf PGPoolConf) (client *PGClient, err error) {
	dsn := conf.GetPGDSN()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	//config.Logger = log15adapter.NewLogger(log.New("module", "pgx"))

	_client, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Postgres!")
	log.Printf("Use Database: \"%s\"\n", conf.DBName)
	return &PGClient{*_client}, err
}

type PGClient struct {
	pgxpool.Pool
}

func (p *PGClient) GetConn(ctx context.Context) (*PGConn, error) {
	conn, err := p.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return &PGConn{*conn}, nil
}

func (p *PGClient) Begin(ctx context.Context) (*PGTx, error) {
	return p.BeginTx(ctx, pgx.TxOptions{})
}

func (p *PGClient) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (*PGTx, error) {
	tx, err := p.Pool.BeginTx(ctx, txOptions)

	if err != nil {
		return nil, err
	}

	return &PGTx{*(tx.(*pgxpool.Tx))}, err
}

type PGConn struct {
	pgxpool.Conn
}

func (c *PGConn) ExecRow(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	commandTag, err := c.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
}

func (c *PGConn) Insert(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return c.ExecRow(ctx, sql, args...)
}

func (c *PGConn) Update(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return c.ExecRow(ctx, sql, args...)
}

func (c *PGConn) Delete(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return c.ExecRow(ctx, sql, args...)
}

func (c *PGConn) FindOne(ctx context.Context, sql string, result interface{}, args ...interface{}) error {
	rows, err := c.Query(ctx, sql, args...)

	if err != nil {
		return err
	}

	return StructScanOne(rows, result)
}

func (c *PGConn) FindAll(ctx context.Context, sql string, result interface{}, args ...interface{}) error {
	rows, err := c.Query(ctx, sql, args...)

	if err != nil {
		return err
	}

	return StructScan(rows, result)
}

func (c *PGConn) FindAllAsMap(ctx context.Context, sql string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := c.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	return PGMapScan(rows)
}

func (c *PGConn) FindAllAsArray(ctx context.Context, sql string, args ...interface{}) ([][]interface{}, error) {
	rows, err := c.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	return PGArrayScan(rows)
}

func (c *PGConn) Count(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	var count int64
	rows, err := c.Query(ctx, sql, args...)

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}

type PGTx struct {
	pgxpool.Tx
}

func (tx *PGTx) RollBackDefer(ctx context.Context) {
	_ = tx.Rollback(ctx)
}

func (tx *PGTx) ExecRow(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	commandTag, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
}

func (tx *PGTx) Insert(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return tx.ExecRow(ctx, sql, args...)
}

func (tx *PGTx) Update(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return tx.ExecRow(ctx, sql, args...)
}

func (tx *PGTx) Delete(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return tx.ExecRow(ctx, sql, args...)
}

func (tx *PGTx) FindOne(ctx context.Context, sql string, result interface{}, args ...interface{}) error {
	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		return err
	}

	return StructScanOne(rows, result)
}

func (tx *PGTx) FindAll(ctx context.Context, sql string, result interface{}, args ...interface{}) error {
	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		return err
	}

	return StructScan(rows, result)
}

func (tx *PGTx) FindAllAsMap(ctx context.Context, sql string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	return PGMapScan(rows)
}

func (tx *PGTx) FindAllAsArray(ctx context.Context, sql string, args ...interface{}) ([][]interface{}, error) {
	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	return PGArrayScan(rows)
}

func (tx *PGTx) Count(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	var count int64
	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}
