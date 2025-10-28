package internal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

var (
	ErrUnsupportedCharset      = errors.New("unsupported charset")
	ErrUnsupportedOutputFormat = errors.New("unsupported output format")
)

func DecodeXML(xmlPath string) ([]*ResultValute, error) {
	file, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("read xml: %w", err)
	}

	reader := bytes.NewReader(file)

	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, charset)
	}

	result := make([]*ResultValute, 0)

	for {
		tkn, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("read XML token: %w", err)
		}

		if start, ok := tkn.(xml.StartElement); ok {
			if start.Name.Local == "Valute" {
				var valuteXML ValuteXML

				if err := decoder.DecodeElement(&valuteXML, &start); err != nil {
					return nil, fmt.Errorf("decode valute element: %w", err)
				}

				numcode, err := strconv.Atoi(valuteXML.NumCode)
				if err != nil {
					numcode = 0
				}

				strValue := strings.ReplaceAll(valuteXML.Value, ",", ".")

				value, err := strconv.ParseFloat(strValue, 64)
				if err != nil {
					value = 0.0
				}

				valute := &ResultValute{numcode, valuteXML.CharCode, value}

				result = append(result, valute)
			}
		}
	}

	return result, nil
}

func EncodeFile(valutes []*ResultValute, outputFormat string, outputPath string) error {
	var (
		encodedData []byte
		err         error
	)

	switch outputFormat {
	case "json":
		encodedData, err = json.MarshalIndent(valutes, "", "    ")
	case "yaml":
		encodedData, err = yaml.Marshal(valutes)
	case "xml":
		encodedData, err = xml.MarshalIndent(valutes, "", "    ")
		if err == nil {
			encodedData = append([]byte(xml.Header), encodedData...)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedOutputFormat, outputFormat)
	}

	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	dir := filepath.Dir(outputPath)

	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("error creating output directory %s: %w", dir, err)
	}

	if err := os.WriteFile(outputPath, encodedData, filePerm); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}
