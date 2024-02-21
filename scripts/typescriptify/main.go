package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	converter := typescriptify.New()
	converter.BackupDir = ""

	converter = converter.
		Add(api.WSData{}).
		Add(api.WSEvent{}).
		Add(api.WSDahuaEvent{}).
		ManageType(time.Time{}, typescriptify.TypeOptions{TSType: "Date", TSTransform: "new Date(__VALUE__)"}).
		ManageType(json.RawMessage{}, typescriptify.TypeOptions{TSType: "Object"})

	err := converter.ConvertToFile(os.Args[1])
	if err != nil {
		panic(err.Error())
	}
}
