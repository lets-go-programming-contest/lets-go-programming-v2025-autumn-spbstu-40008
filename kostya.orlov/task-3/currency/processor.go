package currency

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

func DecodeXML(xmlPath string) ([]*ResultValute, error) {
	file, err := os.ReadFile(xmlPath)

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(file)

	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unknowm charset: %s", charset)
	}

	var valcurs ValCurs

	err = decoder.Decode(&valcurs)

	if err != nil {
		return nil, err
	}

	var result []*ResultValute

	for _, elem := range valcurs.Valutes {
		numcode, err := strconv.Atoi(elem.NumCode)

		if err != nil {
			numcode = 0
		}

		strValue := elem.Value
		strValue = strings.Replace(strValue, ",", ".", -1)

		value, err := strconv.ParseFloat(strValue, 64)

		if err != nil {
			value = 0.0
		}

		valute := &ResultValute{numcode, elem.CharCode, value}

		result = append(result, valute)
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
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	if err != nil {
		return fmt.Errorf("error marshalling data to %s: %w", outputFormat, err)
	}

	dir := filepath.Dir(outputPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating output directory %s: %w", dir, err)
	}

	if err := os.WriteFile(outputPath, encodedData, 0644); err != nil {
		return fmt.Errorf("error writing to file %s: %w", outputPath, err)
	}

	return nil
}
