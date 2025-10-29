package xmlpack

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/mordw1n/task-3/internal/valute"
)

func ReadFile(filePath string) valute.ValuteCurs {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Read XML %q: %w", path, err)
	}
	var valCurs valute.ValuteCurs
	if err := xml.Unmarshal(data, &valCurs); err != nil {
		return valCurs, fmt.Errorf("Unmarshal %q: %w", path, err)
	}
	return valCurs, nil
}