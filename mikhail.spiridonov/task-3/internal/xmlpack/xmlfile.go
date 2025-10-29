package xmlpack

import (
	"os"
	"encoding/xml"

	"github.com/mordw1n/task-3/internal/valute"
)

func ReadFile(filePath string) valute.ValuteCurs {
	data, _ := os.ReadFile(filePath)
	var valCurs valute.ValuteCurs
	xml.Unmarshal(data, &valCurs)
	return valCurs
}