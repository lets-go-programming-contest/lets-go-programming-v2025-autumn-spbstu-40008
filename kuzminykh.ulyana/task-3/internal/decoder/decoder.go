package decoder

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/kuzminykh.ulyana/task-3/internal/models"
	"golang.org/x/text/encoding/charmap"
)

func DecodeFile(filePath string) (*models.ValCurs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	decoder := charmap.Windows1251.NewDecoder()
	decodedData, err := decoder.Bytes(data)
	if err != nil {
		return nil, fmt.Errorf("encoding conversion error: %w", err)
	}

	xmlStr := string(decodedData)
	if idx := strings.Index(xmlStr, "?>"); idx != -1 {
		xmlStr = xmlStr[idx+2:]
	}

	var currencies models.ValCurs
	if err := xml.Unmarshal([]byte(xmlStr), &currencies); err != nil {
		return nil, fmt.Errorf("unmarshaling XML: %w", err)
	}

	return &currencies, nil
}
