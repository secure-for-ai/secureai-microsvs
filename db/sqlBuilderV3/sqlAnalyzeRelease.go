//go:build !debug

package sqlBuilderV3

import (
	"context"

	"github.com/secure-for-ai/secureai-microsvs/db/pgdb"
)

func analyzeQuery(q pgdb.PGQuerier, ctx context.Context, sql string, args ...any) {}