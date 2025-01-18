//go:build debug

package sqlBuilderV3

import (
	"context"
	"crypto"
	"encoding/hex"

	"github.com/secure-for-ai/secureai-microsvs/db/pgdb"
	"github.com/secure-for-ai/secureai-microsvs/log"
	"github.com/secure-for-ai/secureai-microsvs/util"
)

func analyzeQuery(q pgdb.PGQuerier, ctx context.Context, sql string, args ...any) {
	// start ANALYZE SQL by printing the sql and its args
	log.Debug("==============================START ANALYZE SQL==============================")
	log.Debugln(sql)
	log.Debugln(args...)

	// Prepare EXPLAIN ANALYZE query
	sqlName := "explainSql_" + hex.EncodeToString(util.HashString(sql, crypto.SHA256))
	_, err := q.Prepare(ctx, sqlName, "EXPLAIN ANALYZE " + sql)

	if err != nil {
		log.Fatalf("Unable to prepare statement: %v\n", err)
	}

	// Begin a transaction
	tx, err := q.Begin(ctx)
	if err != nil {
		log.Fatalf("Unable to begin transaction: %v\n", err)
	}

	// Execute EXPLAIN ANALYZE with parameters for the prepared statement
	rows, err := tx.Query(ctx, sqlName, args...)
	if err != nil {
		log.Fatalf("Failed to execute EXPLAIN ANALYZE: %v\n", err)
	}
	defer rows.Close()

	// Print the results of EXPLAIN ANALYZE
	log.Debugln("EXPLAIN ANALYZE output:")
	for rows.Next() {
		var output string
		if err := rows.Scan(&output); err != nil {
			log.Fatalf("Failed to scan row: %v\n", err)
		}
		log.Debugln(output)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v\n", err)
	}

	// Rollback the transaction to avoid modifying data
	if err := tx.Rollback(ctx); err != nil {
		log.Fatalf("Failed to rollback transaction: %v\n", err)
	}

	// end ANALYZE SQL
	log.Debugln("===============================END ANALYZE SQL===============================")
}