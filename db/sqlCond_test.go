package db_test

import (
	"github.com/secure-for-ai/secureai-microsvs/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	condNull = db.Expr("")
	cond1    = db.Expr("A < ?", 1)
	cond2    = db.Expr("B = ?", "hello")
	cond3    = db.Expr("C LIKE ?", "username")
)

func TestExpr(t *testing.T) {
	sql1, args, err := db.CondToSQL(cond1)
	assert.NoError(t, err)
	assert.EqualValues(t, "A < ?", sql1)
	assert.EqualValues(t, []interface{}{1}, args)
}

func TestExpr_And(t *testing.T) {
	cond1And2 := cond1.And(cond2)
	sql2, args, err := db.CondToSQL(cond1And2)
	assert.NoError(t, err)
	assert.EqualValues(t, "(A < ?) AND (B = ?)", sql2)
	assert.EqualValues(t, []interface{}{1, "hello"}, args)

	//test And with empty cond
	tmpCond := cond1.And(condNull)
	sql2, args, err = db.CondToSQL(tmpCond)
	assert.NoError(t, err)
	assert.EqualValues(t, cond1, tmpCond)
	assert.EqualValues(t, "A < ?", sql2)
	assert.EqualValues(t, []interface{}{1}, args)

	//test nest And
	cond1And2And3And1And2 := cond1And2.And(cond3, condNull, condNull, cond1And2)
	sql2, args, err = db.CondToSQL(cond1And2And3And1And2)
	assert.NoError(t, err)
	assert.EqualValues(t, "(A < ?) AND (B = ?) AND (C LIKE ?) AND (A < ?) AND (B = ?)", sql2)
	assert.EqualValues(t, []interface{}{1, "hello", "username", 1, "hello"}, args)
}

func TestExpr_Or(t *testing.T) {
	cond1Or2 := cond1.Or(cond2)
	sql2, args, err := db.CondToSQL(cond1Or2)
	assert.NoError(t, err)
	assert.EqualValues(t, "(A < ?) OR (B = ?)", sql2)
	assert.EqualValues(t, []interface{}{1, "hello"}, args)

	//test Or with empty cond
	tmpCond := cond1.Or(condNull)
	sql2, args, err = db.CondToSQL(tmpCond)
	assert.NoError(t, err)
	assert.EqualValues(t, cond1, tmpCond)
	assert.EqualValues(t, "A < ?", sql2)
	assert.EqualValues(t, []interface{}{1}, args)

	//test nest And
	cond1Or2Or3Or1Or2 := cond1Or2.Or(cond3, condNull, condNull, cond1Or2)
	sql2, args, err = db.CondToSQL(cond1Or2Or3Or1Or2)
	assert.NoError(t, err)
	assert.EqualValues(t, "(A < ?) OR (B = ?) OR (C LIKE ?) OR (A < ?) OR (B = ?)", sql2)
	assert.EqualValues(t, []interface{}{1, "hello", "username", 1, "hello"}, args)
}

func TestAnd(t *testing.T) {
	assert.EqualValues(t, condNull, db.NewCond().And().And(condNull).And(condNull, condNull))
	assert.EqualValues(t, condNull, db.And().And(condNull).And(condNull, condNull))
	assert.EqualValues(t, cond1, db.And(cond1))
	assert.EqualValues(t, cond1, db.And(condNull, cond1))
	assert.EqualValues(t, cond1, db.And(cond1, condNull))
	assert.EqualValues(t, cond1, db.And(condNull, cond1, condNull))
	assert.EqualValues(t, cond1, condNull.And(condNull, cond1, condNull).And(condNull, condNull))
}

func TestOr(t *testing.T) {
	assert.EqualValues(t, condNull, db.NewCond().Or().Or(condNull).Or(condNull, condNull))
	assert.EqualValues(t, condNull, db.Or().Or(condNull).Or(condNull, condNull))
	assert.EqualValues(t, cond1, db.Or(cond1))
	assert.EqualValues(t, cond1, db.Or(condNull, cond1))
	assert.EqualValues(t, cond1, db.Or(cond1, condNull))
	assert.EqualValues(t, cond1, db.Or(condNull, cond1, condNull))
	assert.EqualValues(t, cond1, condNull.Or(condNull, cond1, condNull).Or(condNull, condNull))
}

func TestNewCond(t *testing.T) {
	assert.EqualValues(t, condNull, db.NewCond().Or())
}

func TestComplex(t *testing.T) {
	assert.EqualValues(t, condNull,
		db.NewCond().And().And(condNull).And(condNull, condNull).
			Or().Or(condNull).Or(condNull, condNull),
	)
	assert.EqualValues(t, condNull,
		db.And().And(condNull).And(condNull, condNull).
			Or().Or(condNull).Or(condNull, condNull),
	)
	assert.EqualValues(t, condNull,
		condNull.
			And().And(condNull).And(condNull, condNull).
			Or().Or(condNull).Or(condNull, condNull),
	)

	cond1And2 := cond1.And(cond2)
	cond1Or3 := cond1.Or(cond3)

	tmpCond := cond1And2.Or(cond1And2).And(cond1Or3)
	sql, args, err := db.CondToSQL(tmpCond)
	assert.NoError(t, err)
	assert.EqualValues(t, "(((A < ?) AND (B = ?)) OR ((A < ?) AND (B = ?))) AND ((A < ?) OR (C LIKE ?))", sql)
	assert.EqualValues(t, []interface{}{1, "hello", 1, "hello", 1, "username"}, args)
}
