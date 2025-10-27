package format

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/narumov-diyar/task-3/internal/currency"
	"gopkg.in/yaml.v3"
)

var ErrUnsupportedFormat = errors.New("unsupported output format")

const outputDirPerm = 0o755

func writeXML(file *os.File, currencies []currency.Currency) error {
	root := struct {
		XMLName    xml.Name            `xml:"ValCurs"`
		Currencies []currency.Currency `xml:"Valute"`
	}{
		XMLName:    xml.Name{Space: "", Local: "ValCurs"},
		Currencies: currencies,
	}

	data, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal XML: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write XML to file: %w", err)
	}

	return nil
}

func Write(currencies []currency.Currency, outputPath, formatName string) error {
	switch formatName {
	case "json", "yaml", "xml":
	default:
		return fmt.Errorf("%w: %q", ErrUnsupportedFormat, formatName)
	}

	base := strings.TrimSuffix(outputPath, filepath.Ext(outputPath))
	outputPath = base + "." + formatName

	if err := os.MkdirAll(filepath.Dir(outputPath), outputDirPerm); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			return
		}
	}()

	switch formatName {
	case "json":
		enc := json.NewEncoder(file)
		enc.SetIndent("", "\t")

		if err := enc.Encode(currencies); err != nil {
			return fmt.Errorf("encode JSON: %w", err)
		}
	case "yaml":
		if err := yaml.NewEncoder(file).Encode(currencies); err != nil {
			return fmt.Errorf("encode YAML: %w", err)
		}
	case "xml":
		return writeXML(file, currencies)
	}

	return nil
}
