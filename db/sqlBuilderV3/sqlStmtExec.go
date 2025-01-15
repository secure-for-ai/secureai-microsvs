package sqlBuilderV3

import (
	"context"
	"crypto"
	"encoding/hex"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/secure-for-ai/secureai-microsvs/db"
	"github.com/secure-for-ai/secureai-microsvs/db/pgdb"
	"github.com/secure-for-ai/secureai-microsvs/util"
)

var sqlCache = sync.Map{}

func (stmt *Stmt) ExecPG(tx pgdb.PGQuerier, ctx context.Context, result ...interface{}) (int64, error) {
	w := NewWriter()
	defer w.Destroy()
	sql, args, err := stmt.Gen(w, db.SchPG)

	// Get the sql from the cache as the sql is hold by w, which is a reusable buffer
	// pgx need to store sql in its local cache, therefore, we need to make a deep copy
	// of the sql.
	if val, ok := sqlCache.Load(sql); ok {
		sql = val.(string)
	} else {
		sql = strings.Clone(sql)
		sqlCache.Store(sql, sql)
	}

	// there is an error in query generation.
	if err != nil {
		return 0, err
	}

	switch stmt.sqlType {
	case InsertType:
		// Insert Select or Insert one record
		if len(stmt.tableFrom) > 0 || len(stmt.InsertValues) == 1 {
			return tx.ExecRowsAffected(ctx, sql, args...)
		}

		// Insert multiple rows

		// ts := strconv.FormatInt(util.GetNowTimestamp(), 10)
		// nonce, _ := util.GenerateRandomKey(8)
		// sqlName := util.Base64EncodeToString(nonce) + ts

		sqlName := "sql_" + hex.EncodeToString(util.HashString(sql, crypto.SHA256))

		_, err := tx.Prepare(ctx, sqlName, sql)
		if err != nil {
			return 0, err
		}

		batch := &pgx.Batch{}

		bulkArgs := w.BulkArgs()
		rows := len(bulkArgs)
		for _, args := range bulkArgs {
			batch.Queue(sqlName, *args...)
		}

		br := tx.SendBatch(context.Background(), batch)

		var affectedRows int64 = 0
		var errs util.MultiError

		for i := 0; i < rows; i++ {
			tag, err := br.Exec()
			if err != nil {
				errs = append(errs, err)
			} else {
				affectedRows += tag.RowsAffected()
			}
		}

		err = br.Close()
		if err != nil {
			errs = append(errs, err)
		}

		if len(errs) == 0 {
			return affectedRows, nil
		}
		return affectedRows, errs
	case DeleteType, UpdateType:
		return tx.ExecRowsAffected(ctx, sql, args...)
	case SelectType:
		rows, err := tx.Query(ctx, sql, args...)

		if err != nil {
			return 0, err
		}

		// result is not given, so do nothing
		if len(result) == 0 {
			rows.Close()
			return rows.CommandTag().RowsAffected(), err
		}

		res := result[0]
		resValue := util.ReflectValue(res)

		switch resValue.Kind() {
		case reflect.Struct:
			err = pgdb.StructScanOne(rows, result[0])
		case reflect.Slice:
			// if the data type of result[0] is a slice, then pre-allocate
			// the memory up to stmt.LimitN slots in case of resValue.Cap() < stmt.LimitN
			if resValue.Cap() < stmt.LimitN {
				resValue.Set(reflect.MakeSlice(resValue.Type(), 0, stmt.LimitN))
			}

			// handle map scan
			if maps, ok := res.(*[]map[string]interface{}); ok {
				err = pgdb.PGMapScan(rows, maps)
				goto RowClose
			}

			// handle array scan
			if arr, ok := res.(*[][]interface{}); ok {
				err = pgdb.PGArrayScan(rows, arr)
				goto RowClose
			}
			err = pgdb.StructScanSlice(rows, result[0])
		default:
			err = errors.New("not support result data type: " + reflect.TypeOf(result[0]).String())
			rows.Close()
		}

	RowClose:
		return rows.CommandTag().RowsAffected(), err
	default:
		return 0, ErrNotSupportType
	}
}
