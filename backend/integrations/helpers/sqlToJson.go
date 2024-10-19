package helpers

import (
    "database/sql"
    _ "encoding/json"
    "strconv"
    _ "fmt"
    _ "reflect"
    _ "github.com/go-sql-driver/mysql"
)

type SqlToJSONParams struct {
    Rows *sql.Rows
    FieldToNormalize string
    NormalizingFunction func(input float64) (result string)
}

// SQLToJSON takes an SQL result and converts it to a nice JSON form. It also
// handles possibly-null values nicely. See https://stackoverflow.com/a/52572145/265521
func SQLToJSON(params SqlToJSONParams) ([]interface {}, error) {
    rows := params.Rows
    columnTypes, err := rows.ColumnTypes()

    if err != nil {
        return nil, err
    }

    count := len(columnTypes)
    finalRows := []interface{}{};

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

            if z, ok := (scanArgs[i]).(*sql.NullBool); ok  {
                masterData[v.Name()] = z.Bool
                continue;
            }

            if z, ok := (scanArgs[i]).(*sql.NullString); ok  {
                if (params.FieldToNormalize != "" && v.Name() == params.FieldToNormalize) {
                    if s, err := strconv.ParseFloat(z.String, 64); err == nil {
                        z.String = params.NormalizingFunction(s);
                    }
                }
                masterData[v.Name()] = z.String
                continue;
            }

            if z, ok := (scanArgs[i]).(*sql.NullInt64); ok  {
                masterData[v.Name()] = z.Int64
                continue;
            }

            if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok  {
                masterData[v.Name()] = z.Float64
                continue;
            }

            if z, ok := (scanArgs[i]).(*sql.NullInt32); ok  {
                masterData[v.Name()] = z.Int32
                continue;
            }

            masterData[v.Name()] = scanArgs[i]
        }

        finalRows = append(finalRows, masterData)
    }

    return finalRows, nil;
}
