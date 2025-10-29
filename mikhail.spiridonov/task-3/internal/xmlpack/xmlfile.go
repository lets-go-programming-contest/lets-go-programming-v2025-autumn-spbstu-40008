package xmlpack

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/mordw1n/task-3/internal/valute"
	"golang.org/x/net/html/charset"
)

func ReadFile(filePath string) (valute.ValuteCurs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return valute.ValuteCurs{}, fmt.Errorf("Open XML %q: %w", filePath, err)
	}
	defer file.Close()

	var valCurs valute.ValuteCurs
	dcdr := xml.NewDecoder(file)
	dcdr.CharsetReader = charset.NewReaderLabel

	if err := dcdr.Decode(&valCurs); err != nil {
		return valute.ValuteCurs{}, fmt.Errorf("Unmarshal %q: %w", filePath, err)
	}

	return valCurs, nil
}
