package debug

import (
	"encoding/json"
	"os"
)

func Dump(key string, data any) {
	bytes, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile("./debug/"+key+".json", bytes, 0777)
}
