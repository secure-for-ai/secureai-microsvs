package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	//"github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
	//"go.mongodb.org/mongo-driver/mongo"
	//"os"
)

//user=jack password=secret host=pg.example.com port=5432 dbname=mydb
type PGPoolConf struct {
	Host   string `json:"Host"`
	Port   string `json:"Port"`
	DBName string `json:"DBName"`
	User   string `json:"User"`
	PW     string `json:"PW"`
}

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
	return &PGClient{client: _client}, err
}

type PGClient struct {
	client *pgxpool.Pool
}

func (p *PGClient) GetConn(ctx context.Context) (*PGConn, error) {
	conn, err := p.client.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return &PGConn{conn: conn}, nil
}

type PGConn struct {
	conn *pgxpool.Conn
}

func (c *PGConn) Release() {
	c.conn.Release()
}

func (c *PGConn) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return c.conn.Exec(ctx, sql, arguments...)
}

func (c *PGConn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.conn.Query(ctx, sql, args...)
}

func (c *PGConn) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.conn.QueryRow(ctx, sql, args...)
}

func (c *PGConn) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return c.conn.SendBatch(ctx, b)
}

func (c *PGConn) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return c.conn.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (c *PGConn) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.conn.Begin(ctx)
}

func (c *PGConn) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.conn.BeginTx(ctx, txOptions)
}

func (c *PGConn) InsertOne(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	commandTag, err := c.conn.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
}

func (c *PGConn) UpdateOne(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	commandTag, err := c.conn.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
}

func (c *PGConn) DeleteOne(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	commandTag, err := c.conn.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), err
}

func (c *PGConn) FindOne(ctx context.Context, sql string, result interface{}, args ...interface{}) error {
	rows, err := c.conn.Query(ctx, sql, args...)

	if err != nil {
		return err
	}

	return StructScanOne(rows, result)
}

func (c *PGConn) FindAll(ctx context.Context, sql string, result interface{}, args ...interface{}) error {
	rows, err := c.conn.Query(ctx, sql, args...)

	if err != nil {
		return err
	}

	return StructScan(rows, result)
}

func (c *PGConn) FindAllAsMap(ctx context.Context, sql string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := c.conn.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	return PGMapScan(rows)
}

func (c *PGConn) FindAllAsArray(ctx context.Context, sql string, args ...interface{}) ([][]interface{}, error) {
	rows, err := c.conn.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	return PGArrayScan(rows)
}

type PGTx struct {
	t pgx.Tx
}

func (tx *PGTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return tx.t.Begin(ctx)
}

func (tx *PGTx) Commit(ctx context.Context) error {
	return tx.t.Commit(ctx)
}

func (tx *PGTx) Rollback(ctx context.Context) error {
	return tx.t.Rollback(ctx)
}

func (tx *PGTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return tx.t.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (tx *PGTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return tx.t.SendBatch(ctx, b)
}

func (tx *PGTx) LargeObjects() pgx.LargeObjects {
	return tx.t.LargeObjects()
}

func (tx *PGTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return tx.t.Prepare(ctx, name, sql)
}

func (tx *PGTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return tx.t.Exec(ctx, sql, arguments...)
}

func (tx *PGTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return tx.t.Query(ctx, sql, args...)
}

func (tx *PGTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return tx.t.QueryRow(ctx, sql, args...)
}

func (tx *PGTx) Conn() *pgx.Conn {
	return tx.t.Conn()
}

//func (x *XConn) Hello( )  {
//	x.(*pgxpool.Conn).Begin(context.Background())
//}
//func main() {
//	pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
//}
