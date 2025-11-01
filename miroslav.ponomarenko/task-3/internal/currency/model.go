package currency

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Decimal float64

var (
	ErrEmptyNumber        = errors.New("empty number")
	ErrMultipleSeparators = errors.New("multiple decimal separators")
	ErrInvalidNumber      = errors.New("invalid number")
)

func (d *Decimal) UnmarshalText(text []byte) error {
	numStr := strings.TrimSpace(string(text))
	if numStr == "" {
		return ErrEmptyNumber
	}

	numStr = strings.Replace(numStr, ",", ".", 1)

	if strings.Contains(numStr, ",") {
		return fmt.Errorf("%w: %q", ErrMultipleSeparators, text)
	}

	file, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return fmt.Errorf("%w: %q: %w", ErrInvalidNumber, text, err)
	}

	*d = Decimal(file)

	return nil
}

type ExchangeRates struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    Decimal `json:"value"     xml:"Value"`
}
