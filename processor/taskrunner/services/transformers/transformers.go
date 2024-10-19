package transformers

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"encoding/csv"
	"encoding/json"

	"github.com/clbanning/mxj/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const TransformerTypeStrSplit = "str_split"
const TransformerTypeExtractFromJSON = "extract_from_json"
const TransformerTypeCastTo = "cast_to"
const TransformerTypeEncode = "encode_to"

const TransformerTypeCastToString = "string"
const TransformerTypeCastToInteger = "integer"

const TransformerTypeEncodeToXML = "XML"
const TransformerTypeEncodeToCSV = "CSV"

const AsciiTabCode = 9

type Xml struct {
	XMLName string `xml:"xml"`
	Items   []interface{}
}

func SplitString(value string, delimiter string) []string {
	return strings.Split(value, delimiter)
}

func CastStringToInteger(value string) (int64, error) {
	n, err := strconv.ParseFloat(value, 64)

	if err != nil {
		return 0, err
	}

	return int64(n), nil
}

func CastFloatToInteger(value float64) int64 {
	return int64(value)
}

func CastToStringType(value float64) string {
	return fmt.Sprintf("%v", value)
}

func EncodeToCsv(value []byte, delimiter string, columnsOrder []string) ([]byte, error) {
	// Unmarshal JSON data
	var jsonData []map[string]interface{}
	err := json.Unmarshal(value, &jsonData)

	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	var d rune
	if delimiter == "\\t" {
		d = rune(AsciiTabCode)
	} else if len(delimiter) > 0 {
		rd := []rune(delimiter)
		d = rd[0]
	} else {
		d = ','
	}
	writer.Comma = d

	var row []string
	row = append(row, columnsOrder...)
	_ = writer.Write(row)

	for _, object := range jsonData {
		var row []string

		for _, v := range columnsOrder {
			value, ok := object[v]
			if !ok {
				log.Infof("No such column %s", v)

				continue
			}

			var stringValue string
			switch in := value.(type) {
			case float64:
				if math.Mod(in, 1.0) == 0 {
					stringValue = fmt.Sprintf("%d", int(in))
				} else {
					stringValue = fmt.Sprintf("%.2f", in)
				}
			case string:
				stringValue = in
			case nil:
				stringValue = ""
			default:
				converted, ok := value.(string)
				if !ok {
					log.Warnf(`An attempt to cast variable to string, but failed. The var "%v"`, in)
					stringValue = "???"
				} else {
					stringValue = converted
				}
			}

			row = append(row, stringValue)
		}

		_ = writer.Write(row)
	}

	writer.Flush()

	return buffer.Bytes(), nil
}

func EncodeToXml(value []byte) ([]byte, error) {
	// Unmarshal JSON data
	var jsonData []interface{}
	err := json.Unmarshal(value, &jsonData)

	if err != nil {
		return nil, err
	}

	out, err := mxj.AnyXmlIndent(jsonData, "", "  ", "xml")

	return out, err
}

func ExtractFromJsonWithJsonQuery(value []byte, expression string) gjson.Result {
	return gjson.GetBytes(value, expression)
}

/*func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}*/
