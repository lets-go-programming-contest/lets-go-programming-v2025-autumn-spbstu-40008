package xmlpack

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/mordw1n/task-3/internal/valute"
)

func ReadFile(filePath string) (valute.ValuteCurs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return valute.ValuteCurs{}, fmt.Errorf("Read XML %q: %w", filePath, err)
	}

	var valCurs valute.ValuteCurs
	if err := xml.Unmarshal(data, &valCurs); err != nil {
		return valute.ValuteCurs{}, fmt.Errorf("Unmarshal %q: %w", filePath, err)
	}

	return valCurs, nil
}