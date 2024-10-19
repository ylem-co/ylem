package tableau

import (
	"fmt"
	"reflect"
)

func Rows(columnOrder []string, rows []map[string]interface{}) [][]interface{} {
	result := make([][]interface{}, len(rows))
	for i, row := range rows {
		rrow := make([]interface{}, len(columnOrder))
		for k, col := range columnOrder {
			rrow[k] = row[col]
		}

		result[i] = rrow
	}

	return result
}

func Columns(columnOrder []string, data []map[string]interface{}) ([]Column, error) {
	result := make([]Column, len(columnOrder))
	if len(data) == 0 {
		return result, nil
	}

	for i, col := range columnOrder {
		tp := typeStr(data[0][col])
		if tp == "" {
			return result, fmt.Errorf("unable to determine type for column %s", col)
		}
		result[i] = Column{
			Name: col,
			Type: tp,
		}
	}

	return result, nil
}

func typeStr(val interface{}) string {

	switch reflect.ValueOf(val).Kind() {
	case reflect.Float32:
	case reflect.Float64:
		return "double"
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Uintptr:
		return "int"
	case reflect.Bool:
		return "bool"
		// case time.Time:
		// 	return "datetime"

	}

	return "varchar"
}
