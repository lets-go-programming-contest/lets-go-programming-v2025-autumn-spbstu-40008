package decoder

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/task-3/internal/config"
	"github.com/task-3/internal/structures"
	"golang.org/x/text/encoding/charmap"
)

var ErrUnsupportedCharset = errors.New("unsupported charset")

func DecodeXML(cfg config.File) (structures.ValCurs, error) {
	xmlFile, err := os.Open(cfg.Input)
	if err != nil {
		return structures.ValCurs{}, fmt.Errorf("couldn't open %s: %w", cfg.Input, err)
	}

	defer func() {
		if err := xmlFile.Close(); err != nil {
			return
		}
	}()

	var val structures.ValCurs

	decoder := xml.NewDecoder(xmlFile)

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, ErrUnsupportedCharset
	}

	err = decoder.Decode(&val)
	if err != nil {
		return structures.ValCurs{}, fmt.Errorf("couldn't decode XML: %w", err)
	}

	return val, nil
}
