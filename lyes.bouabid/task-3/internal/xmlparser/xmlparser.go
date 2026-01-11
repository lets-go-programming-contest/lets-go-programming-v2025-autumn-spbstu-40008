package xmlparser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ReadXML(path string, result interface{}) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("we cant read xml file: %w", err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(fileData))
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(result)
	if err != nil {
		return fmt.Errorf("parse xml file in failed: %w", err)
	}

	return nil
}
