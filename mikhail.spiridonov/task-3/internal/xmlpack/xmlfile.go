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
		return valute.ValuteCurs{}, fmt.Errorf("open XML %q: %w", filePath, err)
	}

	defer func() {
		if fileErr := file.Close(); fileErr != nil {
			panic(fmt.Errorf("file closed %q: %w", filePath, fileErr))
		}
	}()

	var valCurs valute.ValuteCurs

	dcdr := xml.NewDecoder(file)
	dcdr.CharsetReader = charset.NewReaderLabel

	if err := dcdr.Decode(&valCurs); err != nil {
		return valute.ValuteCurs{}, fmt.Errorf("unmarshal %q: %w", filePath, err)
	}

	return valCurs, nil
}
