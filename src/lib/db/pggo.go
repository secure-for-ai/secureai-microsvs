package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"os"
)

func main() {
	pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
}
