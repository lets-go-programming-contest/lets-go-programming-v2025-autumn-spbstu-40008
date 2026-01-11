package converter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrEmptyNumber        = errors.New("empty number")
	ErrMultipleSeparators = errors.New("multiple decimal separators")
)

type Rates struct {
	Data []Currency `xml:"Valute"`
}

type (
	FloatforCur float64
	Currency    struct {
		NumCode  int         `json:"num_code"  xml:"NumCode"`
		CharCode string      `json:"char_code" xml:"CharCode"`
		Value    FloatforCur `json:"value"     xml:"Value"`
	}
)

func (cf *FloatforCur) UnmarshalText(text []byte) error {
	input := strings.TrimSpace(string(text))
	if input == "" {
		return ErrEmptyNumber
	}

	normalized := strings.Replace(input, ",", ".", 1)
	if strings.Contains(normalized, ",") {
		return ErrMultipleSeparators
	}

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return fmt.Errorf("the number is invalid %q: %w", text, err)
	}

	*cf = FloatforCur(value)

	return nil
}
