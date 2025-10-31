package data

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/task-3/config"
	"github.com/task-3/internal/structures"
	"golang.org/x/text/encoding/charmap"
)

var ErrUnsupportedCharset = errors.New("unsupported charset")

func DecodeXML(cfg config.File) (structures.ReadingXML, error) {
	xmlFile, err := os.Open(cfg.Input)
	if err != nil {
		return structures.ReadingXML{}, fmt.Errorf("failed to open XML input file %s: %w", cfg.Input, err)
	}

	defer func() {
		_ = xmlFile.Close()
	}()

	var xmlData structures.ReadingXML

	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, ErrUnsupportedCharset
	}

	err = decoder.Decode(&xmlData)
	if err != nil {
		return structures.ReadingXML{}, fmt.Errorf("failed to decode XML from file %s: %w", cfg.Input, err)
	}

	return xmlData, nil
} //ы

func ProcessAndSortCurrencies(xmlData structures.ReadingXML) []structures.Currency {
	processed := make([]structures.Currency, 0, len(xmlData.Information))

	for _, item := range xmlData.Information {
		numCode, err := strconv.Atoi(strings.TrimSpace(item.NumCodeStr))

		if err != nil {
			item.NumCode = 0
		} else {
			item.NumCode = numCode
		}

		stringValue := strings.ReplaceAll(item.ValueStr, ",", ".")

		value, err := strconv.ParseFloat(stringValue, 64)
		if err != nil {
			value = 0.0
		}

		nominal, err := strconv.Atoi(strings.TrimSpace(item.NominalStr))

		if err != nil || nominal == 0 {
			nominal = 1
		}

		item.Value = value / float64(nominal)
		item.CharCode = strings.TrimSpace(item.CharCode)

		processed = append(processed, item)
	}

	sort.Slice(processed, func(i, j int) bool {
		return processed[i].Value > processed[j].Value
	})

	return processed
}

func CreateAndWriteJSON(filename string, data []structures.Currency) error {
	dirPath := filepath.Dir(filename)

	const DirPerm = 0o755

	if err := os.MkdirAll(dirPath, DirPerm); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", dirPath, err)
	}

	const FilePerm = 0o644

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePerm)
	if err != nil {
		return fmt.Errorf("failed to open or create output file %s: %w", filename, err)
	}

	defer func() {
		_ = file.Close()
	}()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write JSON data to file: %w", err)
	}

	return nil
}
