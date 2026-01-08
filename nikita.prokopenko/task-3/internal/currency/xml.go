package currency

import (
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"strings"
)

type valCurs struct {
	Valutes []valute `xml:"Valute"`
}

type valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Value    string `xml:"Value"`
}

func DecodeXMLFile(path string) ([]Currency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return DecodeXML(f)
}

func DecodeXML(r io.Reader) ([]Currency, error) {
	var vc valCurs
	if err := xml.NewDecoder(r).Decode(&vc); err != nil {
		return nil, err
	}

	res := make([]Currency, 0, len(vc.Valutes))
	for _, v := range vc.Valutes {
		num := 0
		if s := strings.TrimSpace(v.NumCode); s != "" {
			if n, err := strconv.Atoi(s); err == nil {
				num = n
			}
		}

		valStr := strings.ReplaceAll(
			strings.TrimSpace(v.Value),
			",",
			".",
		)
		valStr = strings.ReplaceAll(valStr, " ", "")
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, err
		}

		nom := 1.0
		if s := strings.TrimSpace(v.Nominal); s != "" {
			s = strings.ReplaceAll(s, ",", ".")
			if n, err := strconv.ParseFloat(s, 64); err == nil && n > 0 {
				nom = n
			}
		}
		if nom != 1 {
			val /= nom
		}

		res = append(res, Currency{
			NumCode:  num,
			CharCode: strings.TrimSpace(v.CharCode),
			Value:    val,
		})
	}
	return res, nil
}
