package helpers

import (
	"fmt"
	"encoding/json"
)

func DumpFormatted(val interface{}) {
	str, _ := json.MarshalIndent(val, "", "    ")
	fmt.Println(string(str))
}
