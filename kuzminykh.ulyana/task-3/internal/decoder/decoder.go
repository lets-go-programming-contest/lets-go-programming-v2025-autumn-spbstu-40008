package decoder

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/kuzminykh.ulyana/task-3/internal/models"
)

func DecodeFile(filePath string) (*models.ValCurs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	content := string(data)
	content = strings.Replace(content, "windows-1251", "UTF-8", 1)

	var currencies models.ValCurs
	if err := xml.Unmarshal([]byte(content), &currencies); err != nil {
		return nil, fmt.Errorf("unmarshaling XML: %w", err)
	}

	return &currencies, nil
}
