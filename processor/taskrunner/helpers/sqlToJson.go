package helpers

import (
	"database/sql"
	_ "encoding/json"
	_ "fmt"
	_ "reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type SqlToJSONParams struct {
	Rows                *sql.Rows
	FieldToNormalize    string
	NormalizingFunction func(input float64) (result string)
}

// SQLToJSON takes an SQL result and converts it to a nice JSON form. It also
// handles possibly-null values nicely. See https://stackoverflow.com/a/52572145/265521
func SQLToJSON(params SqlToJSONParams) ([]interface{}, error) {
	rows := params.Rows
	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return nil, err
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {

		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {

			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
			case "INT":
				scanArgs[i] = new(sql.NullInt32)
			case "DECIMAL", "FLOAT":
				scanArgs[i] = new(sql.NullFloat64)
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)

		if err != nil {
			return nil, err
		}

		masterData := map[string]interface{}{}

		for i, v := range columnTypes {

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				var val *bool
				if z.Valid {
					val = &z.Bool
				}
				masterData[v.Name()] = val

				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				var val *string
				if z.Valid {
					if params.FieldToNormalize != "" && v.Name() == params.FieldToNormalize {
						if s, err := strconv.ParseFloat(z.String, 64); err == nil {
							z.String = params.NormalizingFunction(s)
						}
					}

					val = &z.String
				}
				masterData[v.Name()] = val

				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				var val *int64
				if z.Valid {
					val = &z.Int64
				}
				masterData[v.Name()] = val

				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				var val *float64
				if z.Valid {
					val = &z.Float64
				}
				masterData[v.Name()] = val

				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				var val *int32
				if z.Valid {
					val = &z.Int32
				}
				masterData[v.Name()] = val

				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}

	return finalRows, nil
}
