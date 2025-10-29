package xmlpack

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/mordw1n/task-3/internal/valute"
	"golang.org/x/net/html/charset"
)

func ReadFile(filePath string) (valute.ValuteCurs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return valute.ValuteCurs{}, fmt.Errorf("Read XML %q: %w", filePath, err)
	}

	var valCurs valute.ValuteCurs
	dcdr := xml.NewDecoder(data)
	dcdr.CharsetReader = charset.NewReaderLabel

	if err := dcdr.Decode(&valCurs); err != nil {
		return valCurs, fmt.Errorf("Decode to %q: %w", filePath, err)
	}

	if err := xml.Unmarshal(data, &valCurs); err != nil {
		return valute.ValuteCurs{}, fmt.Errorf("Unmarshal %q: %w", filePath, err)
	}

	return valCurs, nil
}