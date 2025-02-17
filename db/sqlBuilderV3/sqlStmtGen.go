package sqlBuilderV3

import (
	"strconv"
	"strings"

	"github.com/secure-for-ai/secureai-microsvs/db"
)

func (stmt *Stmt) Gen(w *Writer, schema ...db.Schema) (string, []any, error) {
	var err error
	w.Reset()

	switch stmt.sqlType {
	case InsertType:
		err = stmt.insertWriteTo(w)
	case DeleteType:
		err = stmt.deleteWriteTo(w)
	case UpdateType:
		err = stmt.updateWriteTo(w)
	case SelectType:
		err = stmt.selectWriteTo(w)
	}

	sql := strings.Clone(w.String())
	var index, i int

	index = strings.Index(sql, db.Para)

	// nothing need to be replaced
	if index < 0 {
		return sql, w.args, err
	}

	//reset memory of the writer
	w.stringWriter.Reset()
	w.Grow(len(sql))

	start := 0
	sepLen := len(db.Para)

	pgFunc := func() {
		w.WriteString(sql[start : start+index])
		w.WriteByte('$')
		w.WriteString(strconv.Itoa(i + 1))
	}

	myFunc := func() {
		w.WriteString(sql[start : start+index])
		w.WriteByte('?')
	}

	var callback = myFunc

	if len(schema) > 0 {
		switch schema[0] {
		case db.SchPG:
			callback = pgFunc
		case db.SchMYSQL:
			w.Grow(len(sql) - len(w.args))
		}
	}

	for i = 0; ; i++ {
		if index == -1 {
			w.WriteString(sql[start:])
			break
		}
		callback()
		start = start + index + sepLen
		index = strings.Index(sql[start:], db.Para)
	}

	return w.String(), w.args, err
}

func (stmt *Stmt) WriteTo(w *Writer) error {
	switch stmt.sqlType {
	case InsertType:
		return stmt.insertWriteTo(w)
	case DeleteType:
		return stmt.deleteWriteTo(w)
	case UpdateType:
		return stmt.updateWriteTo(w)
	case SelectType:
		return stmt.selectWriteTo(w)
	}
	return ErrNotSupportType
}

func (stmt *Stmt) insertSelectWriteTo(w *Writer) error {
	w.WriteString("INSERT INTO ")
	w.WriteString(stmt.tableInto)

	if len(stmt.InsertCols) > 0 {
		w.WriteString(" (")
		w.Join(stmt.InsertCols, ',')
		w.WriteString(") ")
	} else {
		w.WriteByte(' ')
	}

	if s, ok := stmt.tableFrom[0].(*fromStmt); ok {
		s.writeTo(w)
		return nil
	}

	return stmt.selectWriteTo(w)
}

func (stmt *Stmt) insertWriteTo(w *Writer) error {
	if len(stmt.tableInto) <= 0 {
		return ErrNoTableName
	}

	// Insert Select
	if stmt.tableInto != "" && len(stmt.tableFrom) > 0 {
		return stmt.insertSelectWriteTo(w)
	}

	w.WriteString("INSERT INTO ")
	w.WriteString(stmt.tableInto)

	if len(stmt.InsertCols) > 0 {
		w.WriteString(" (")
		w.Join(stmt.InsertCols, ',')
		w.WriteString(") VALUES (")
	} else {
		w.WriteString(" VALUES (")
	}

	switch rowsLen := len(stmt.InsertValues); rowsLen {
	case 0:
		return ErrNoValueToInsert
	case 1:
		values := stmt.InsertValues[0]
		valuesLen := len(*values)

		for i, value := range *values {
			w.WriteString(value.String())
			w.Append(value.args...)
			if i != valuesLen-1 {
				w.WriteByte(',')
			}
		}
	default:
		// write the first row including sql concat
		values := stmt.InsertValues[0]
		valuesLen := len(*values)
		args := getArgs()

		for i, value := range *values {
			w.WriteString(value.String())
			*args = append(*args, value.args...)
			if i != valuesLen-1 {
				w.WriteByte(',')
			}
		}
		w.AppendBulk(args)

		// write the rest rows
		for _, values := range stmt.InsertValues[1:] {
			args := getArgs()
			for _, value := range *values {
				*args = append(*args, value.args...)
			}
			w.AppendBulk(args)
		}
	}

	w.WriteByte(')')

	return nil
}

func (stmt *Stmt) deleteWriteTo(w *Writer) error {
	if len(stmt.tableFrom) <= 0 {
		return ErrNoTableName
	}

	w.WriteString("DELETE FROM ")
	stmt.tableFrom[0].writeTo(w)

	if stmt.where.IsValid() {
		w.WriteString(" WHERE ")
		stmt.where.WriteTo(w)
	}

	return nil
}

func (stmt *Stmt) updateWriteTo(w *Writer) error {
	if len(stmt.tableFrom) <= 0 {
		return ErrNoTableName
	}

	w.WriteString("UPDATE ")
	stmt.tableFrom[0].writeTo(w)
	w.WriteString(" SET ")
	stmt.SetCols.WriteTo(w)

	if stmt.where.IsValid() {
		w.WriteString(" WHERE ")
		stmt.where.WriteTo(w)
	}

	return nil
}

func (stmt *Stmt) selectWriteTo(w *Writer) error {
	if len(stmt.tableFrom) <= 0 {
		return ErrNoTableName
	}

	w.WriteString("SELECT ")

	if len(stmt.SelectCols) > 0 {
		w.Join(stmt.SelectCols, ',')
	} else {
		w.WriteByte('*')
	}

	w.WriteString(" FROM ")

	for i, from := range stmt.tableFrom {
		from.writeTo(w)
		if i != len(stmt.tableFrom)-1 {
			w.WriteByte(',')
		}
	}

	if stmt.where.IsValid() {
		w.WriteString(" WHERE ")
		stmt.where.WriteTo(w)
	}

	if stmt.GroupByStr.Len() > 0 {
		w.WriteString(" GROUP BY ")
		w.Write(stmt.GroupByStr.Bytes())
	}

	if stmt.having.IsValid() {
		w.WriteString(" HAVING ")
		stmt.having.WriteTo(w)
	}

	if stmt.OrderByStr.Len() > 0 {
		w.WriteString(" ORDER BY ")
		w.Write(stmt.OrderByStr.Bytes())
	}

	if stmt.LimitN < 0 || stmt.Offset < 0 {
		return ErrInvalidLimitation
	} else if stmt.LimitN > 0 {
		w.WriteString(" LIMIT ")
		w.WriteString(strconv.Itoa(stmt.LimitN))
		if stmt.Offset != 0 {
			w.WriteString(" OFFSET ")
			w.WriteString(strconv.Itoa(stmt.Offset))
		}
	}

	return nil
}
