package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"log"

	//"github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PGQuerier interface {
	pgxtype.Querier
	Prepare(ctx context.Context, name string, sql string) (sd *pgconn.StatementDescription, err error)
	ExecRowsAffected(ctx context.Context, sql string, args ...interface{}) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Deallocate(ctx context.Context, name string) error
}

//user=jack password=secret host=pg.example.com port=5432 dbname=mydb
type PGPoolConf struct {
	Host    string `json:"Host"`
	Port    string `json:"Port"`
	DBName  string `json:"DBName"`
	User    string `json:"User"`
	PW      string `json:"PW"`
	Verbose bool   `json:"Verbose"`
}

//type PGMultiError struct{
//	util.MultiError
//}
//
//type PGError string
//
//func (e PGError) Error() string { return string(e) }

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
	if conf.Verbose {
		log.Println("Connected to Postgres!")
		log.Printf("Use Database: \"%s\"\n", conf.DBName)
	}
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

func (p *PGClient) Begin(ctx context.Context) (*PGClientTx, error) {
	return p.BeginTx(ctx, pgx.TxOptions{})
}

func (p *PGClient) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (*PGClientTx, error) {
	tx, err := p.Pool.BeginTx(ctx, txOptions)

	if err != nil {
		return nil, err
	}

	return &PGClientTx{*(tx.(*pgxpool.Tx))}, err
}

type PGConn struct {
	pgxpool.Conn
}

func (c *PGConn) ExecRowsAffected(ctx context.Context, sql string, args ...interface{}) (int64, error) {

	commandTag, err := c.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
}

func (c *PGConn) Insert(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return c.ExecRowsAffected(ctx, sql, args...)
}

func (c *PGConn) Update(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return c.ExecRowsAffected(ctx, sql, args...)
}

func (c *PGConn) Delete(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return c.ExecRowsAffected(ctx, sql, args...)
}

func (c *PGConn) FindOne(ctx context.Context, sql string, result interface{}, args ...interface{}) error {
	rows, err := c.Query(ctx, sql, args...)

	if err != nil {
		return err
	}

	return StructScanOne(rows, result)
}

func (c *PGConn) FindAll(ctx context.Context, sql string, result interface{}, args ...interface{}) (int64, error) {
	rows, err := c.Query(ctx, sql, args...)

	if err != nil {
		return 0, err
	}

	err = StructScanSlice(rows, result)

	if err != nil {
		return 0, err
	}

	return rows.CommandTag().RowsAffected(), err
}

func (c *PGConn) FindAllAsMap(ctx context.Context, sql string, result *[]map[string]interface{}, args ...interface{}) (int64, error) {
	rows, err := c.Query(ctx, sql, args...)

	if err != nil {
		return 0, err
	}

	err = PGMapScan(rows, result)

	if err != nil {
		return 0, err
	}

	return rows.CommandTag().RowsAffected(), err
}

func (c *PGConn) FindAllAsArray(ctx context.Context, sql string, result *[][]interface{}, args ...interface{}) (int64, error) {
	rows, err := c.Query(ctx, sql, args...)

	if err != nil {
		return 0, err
	}

	err = PGArrayScan(rows, result)

	if err != nil {
		return 0, err
	}

	return rows.CommandTag().RowsAffected(), err
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

func (c *PGConn) Prepare(ctx context.Context, name string, sql string) (sd *pgconn.StatementDescription, err error) {
	return c.Conn.Conn().Prepare(ctx, name, sql)
}

func (c *PGConn) Deallocate(ctx context.Context, name string) error {
	return c.Conn.Conn().Deallocate(ctx, name)
}

type PGClientTx struct {
	pgxpool.Tx
}

func (tx *PGClientTx) RollBackDefer(ctx context.Context) {
	_ = tx.Rollback(ctx)
}

func (tx *PGClientTx) ExecRowsAffected(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	commandTag, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
}

func (tx *PGClientTx) Insert(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return tx.ExecRowsAffected(ctx, sql, args...)
}

func (tx *PGClientTx) Update(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return tx.ExecRowsAffected(ctx, sql, args...)
}

func (tx *PGClientTx) Delete(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	return tx.ExecRowsAffected(ctx, sql, args...)
}

func (tx *PGClientTx) FindOne(ctx context.Context, sql string, result interface{}, args ...interface{}) error {
	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		return err
	}

	return StructScanOne(rows, result)
}

func (tx *PGClientTx) FindAll(ctx context.Context, sql string, result interface{}, args ...interface{}) (int64, error) {
	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		return 0, err
	}

	err = StructScanSlice(rows, result)

	if err != nil {
		return 0, err
	}

	return rows.CommandTag().RowsAffected(), err
}

func (tx *PGClientTx) FindAllAsMap(ctx context.Context, sql string, result *[]map[string]interface{}, args ...interface{}) (int64, error) {
	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		return 0, err
	}

	err = PGMapScan(rows, result)

	if err != nil {
		return 0, err
	}

	return rows.CommandTag().RowsAffected(), err
}

func (tx *PGClientTx) FindAllAsArray(ctx context.Context, sql string, result *[][]interface{}, args ...interface{}) (int64, error) {
	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		return 0, err
	}

	err = PGArrayScan(rows, result)

	if err != nil {
		return 0, err
	}

	return rows.CommandTag().RowsAffected(), err
}

func (tx *PGClientTx) Count(ctx context.Context, sql string, args ...interface{}) (int64, error) {
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

func (tx *PGClientTx) Deallocate(ctx context.Context, name string) error {
	return tx.Conn().Deallocate(ctx, name)
}
