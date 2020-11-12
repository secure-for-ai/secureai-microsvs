package db

import (
	"errors"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"reflect"
)

var ErrFindNil = PGError("pg: row not found")

func getFieldMap(t reflect.Type, rowFields []pgproto3.FieldDescription) []int {
	tNumField := t.NumField()
	// map the field name to field index for a struct
	// Suppose we have struct
	//
	// type User struct {
	// Uid         int64     `json:"uid" pg:"uid"`
	// Username    string    `json:"username" pg:"username"`
	// Password    string    `json:"password" pg:"password"`
	//
	// Then, tFieldNameIndexMap is map[string]int{"password":2, "uid":0, "username":1}
	tFieldNameIndexMap := make(map[string]int, tNumField)
	// Suppose the row field is
	// [
	//    {
	//         "Name":"uid",
	//         "TableOID":16385,
	//         "TableAttributeNumber":1,
	//         "DataTypeOID":20,
	//         "DataTypeSize":8,
	//         "TypeModifier":-1,
	//         "Format":1},
	//    {
	//         "Name":"username",
	//         "TableOID":16385,
	//         "TableAttributeNumber":2,
	//         "DataTypeOID":1043,
	//         "DataTypeSize":-1,
	//         "TypeModifier":104,
	//         "Format":0},
	//    {
	//         "Name":"password",
	//         "TableOID":16385,
	//         "TableAttributeNumber":3,
	//         "DataTypeOID":1043,
	//         "DataTypeSize":-1,
	//         "TypeModifier":104,
	//         "Format":0}
	// ]
	// Then tRowFieldMap = []int{0, 1, 2}
	tRowFieldMap := make([]int, len(rowFields))

	for i := 0; i < tNumField; i++ {
		_field := t.Field(i)
		_tag := _field.Tag.Get("pg")
		if _tag != "" {
			tFieldNameIndexMap[_tag] = i
		} else {
			tFieldNameIndexMap[_field.Name] = i
		}
	}

	for i := range rowFields {
		tRowFieldMap[i] = tFieldNameIndexMap[string(rowFields[i].Name)]
	}
	return tRowFieldMap
}

func StructScanOne(rows pgx.Rows, dest interface{}) error {
	defer rows.Close()

	// get dest ptr
	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	// get real value
	baseValue := reflect.Indirect(value)
	baseType := baseValue.Type()

	fields := rows.FieldDescriptions()
	baseFieldMap := getFieldMap(baseValue.Type(), fields)
	fieldsLen := len(fields)

	vp := reflect.New(baseType)
	v := reflect.Indirect(vp)

	if rows.Next() {
		args := make([]interface{}, fieldsLen)

		for i := range fields {
			args[i] = v.Field(baseFieldMap[i]).Addr().Interface()
		}
		err := rows.Scan(args...)
		if err != nil {
			return err
		}
	} else {
		return ErrFindNil
	}

	baseValue.Set(v)

	return nil
}

func StructScan(rows pgx.Rows, dest interface{}) error {
	var v, vp reflect.Value
	defer rows.Close()

	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	direct := reflect.Indirect(value)

	slice := direct.Type()
	if slice.Kind() != reflect.Slice {
		return errors.New("must pass a pointer to a slice")
	}

	base := slice.Elem()
	if base.Kind() != reflect.Struct {
		panic("reflect: Field of non-struct type " + base.String())
	}

	fields := rows.FieldDescriptions()
	baseFieldMap := getFieldMap(base, fields)
	fieldsLen := len(fields)

	for rows.Next() {
		// allocate new value
		vp = reflect.New(base)
		v = reflect.Indirect(vp)
		args := make([]interface{}, fieldsLen)

		for i := range fields {
			args[i] = v.Field(baseFieldMap[i]).Addr().Interface()
		}
		err := rows.Scan(args...)
		if err != nil {
			return err
		}

		direct.Set(reflect.Append(direct, v))
	}

	return nil
}

func PGMapScan(rows pgx.Rows) ([]map[string]interface{}, error) {

	var maps []map[string]interface{}
	defer rows.Close()

	var m map[string]interface{}
	for rows.Next() {

		v, err := rows.Values()
		if err != nil {
			return nil, err
		}

		fields := rows.FieldDescriptions()

		m = make(map[string]interface{})
		for i := range fields {
			m[string(fields[i].Name)] = v[i]
		}

		maps = append(maps, m)
	}

	return maps, nil
}

func PGArrayScan(rows pgx.Rows) ([][]interface{}, error) {
	var values [][]interface{}
	defer rows.Close()

	for rows.Next() {
		v, err := rows.Values()
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}

	return values, nil
}
