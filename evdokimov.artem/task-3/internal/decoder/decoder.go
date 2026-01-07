package decoder

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/task-3/internal/config"
	"github.com/task-3/internal/structures"
	"golang.org/x/text/encoding/charmap"
)

var (
	ErrUnsupportedCharset = fmt.Errorf("unsupported charset")
)

func DecodeXML(cfg config.File) (structures.ValCurs, error) {
	file, err := os.Open(cfg.Input)
	if err != nil {
		return structures.ValCurs{}, fmt.Errorf("open %s: %w", cfg.Input, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close file: %v\n", closeErr)
		}
	}()

	var result structures.ValCurs

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charsetReader

	if err := decoder.Decode(&result); err != nil {
		return structures.ValCurs{}, fmt.Errorf("decode XML: %w", err)
	}

	return result, nil
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, charset)
	}
}
