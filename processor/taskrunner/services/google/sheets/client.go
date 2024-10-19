package sheets

import (
	"context"
	"fmt"

	messaging "github.com/ylem-co/shared-messaging"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Client interface {
	WriteData(spreadsheetId string, sheetId int64, mode string, sqlQueryColumnOrder []string, rows []map[string]interface{}, writeHeader bool) error
}

type client struct {
	srv *sheets.Service
}

func (c *client) Init(
	ctx context.Context,
	credentials []byte,
) error {
	s, err := sheets.NewService(
		ctx,
		option.WithCredentialsJSON(credentials),
	)

	c.srv = s

	return err
}

func (c *client) WriteData(
	spreadsheetId string,
	sheetId int64,
	mode string,
	sqlQueryColumnOrder []string,
	rows []map[string]interface{},
	writeHeader bool,
) error {
	spreadsheet, err := c.srv.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		log.Error(err)
		return err
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.SheetId == sheetId {
			values := c.rowsToValues(sqlQueryColumnOrder, rows, writeHeader)
			return c.writeData(spreadsheetId, sheet, mode, values, writeHeader)
		}
	}

	return fmt.Errorf("sheet with id %d not found in spreadsheet %s", sheetId, spreadsheetId)
}

func (c *client) writeData(spreadsheetId string, sheet *sheets.Sheet, mode string, values [][]interface{}, writeHeader bool) error {
	if len(values) == 0 {
		return nil
	}

	switch mode {
	case messaging.GoogleSheetsModeAppend:
		return c.appendData(spreadsheetId, sheet, values, writeHeader)

	case messaging.GoogleSheetsModeOverwrite:
		return c.overwriteData(spreadsheetId, sheet, values)
	}

	return fmt.Errorf("unknown google sheets mode: %s", mode)
}

func (c *client) overwriteData(spreadsheetId string, sheet *sheets.Sheet, values [][]interface{}) error {
	rng := fmt.Sprintf("'%s'", sheet.Properties.Title)
	_, err := c.srv.Spreadsheets.Values.Clear(spreadsheetId, rng, &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		return err
	}

	rng = fmt.Sprintf("'%s'!R1C1:R%dC%d", sheet.Properties.Title, len(values), len(values[0]))
	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values,
	}

	_, err = c.srv.Spreadsheets.Values.Update(spreadsheetId, rng, valueRange).ValueInputOption("RAW").Do()
	return err
}

func (c *client) appendData(spreadsheetId string, sheet *sheets.Sheet, values [][]interface{}, writeHeader bool) error {
	startIndex := 0
	if writeHeader {
		startIndex = 1
		rng := fmt.Sprintf("'%s'!R1C1:R1C%d", sheet.Properties.Title, len(values[0]))
		valueRange := &sheets.ValueRange{
			MajorDimension: "ROWS",
			Values:         [][]interface{}{values[0]},
		}

		_, err := c.srv.Spreadsheets.Values.Update(spreadsheetId, rng, valueRange).ValueInputOption("RAW").Do()
		if err != nil {
			return err
		}
	}

	rng := fmt.Sprintf("'%s'", sheet.Properties.Title)
	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values[startIndex:],
	}

	_, err := c.srv.Spreadsheets.Values.Append(spreadsheetId, rng, valueRange).ValueInputOption("RAW").Do()
	return err
}

func (c *client) rowsToValues(sqlQueryColumnOrder []string, rows []map[string]interface{}, writeHeader bool) [][]interface{} {
	var header []string
	headersMatch := true
	if len(rows) > 0 {
		r := rows[0]
		for _, v := range sqlQueryColumnOrder {
			if _, ok := r[v]; !ok {
				headersMatch = false
				break
			}
		}

		headersMatch = headersMatch && len(sqlQueryColumnOrder) == len(r)
	}

	if len(rows) == 0 || headersMatch {
		header = sqlQueryColumnOrder
	} else {
		header = []string{}
		for k := range rows[0] {
			header = append(header, k)
		}
	}

	rowCount := len(rows)
	if writeHeader {
		rowCount++
	}

	values := make([][]interface{}, rowCount)
	if writeHeader {
		values[0] = make([]interface{}, len(header))
		for k, v := range header {
			values[0][k] = v
		}
	}

	for k, row := range rows {
		i := k
		if writeHeader {
			i++
		}

		values[i] = make([]interface{}, len(row))
		for z, fieldName := range header {
			values[i][z] = row[fieldName]
		}
	}

	return values
}

func NewClient(ctx context.Context, credentials []byte) (Client, error) {
	c := &client{}
	err := c.Init(ctx, credentials)
	if err != nil {
		return nil, err
	}

	return c, nil
}
