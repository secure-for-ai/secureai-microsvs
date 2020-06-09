package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"os"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
}
