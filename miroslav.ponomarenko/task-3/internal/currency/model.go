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
	s := strings.TrimSpace(string(text))
	if s == "" {
		return ErrEmptyNumber
	}
	s = strings.Replace(s, ",", ".", 1)
	if strings.Contains(s, ",") {
		return fmt.Errorf("%w: %q", ErrMultipleSeparators, text)
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("%w: %q: %w", ErrInvalidNumber, text, err)
	}
	*d = Decimal(f)
	return nil
}

type ExchangeRates struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int     `xml:"NumCode"  json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    Decimal `xml:"Value" json:"value"`
}
