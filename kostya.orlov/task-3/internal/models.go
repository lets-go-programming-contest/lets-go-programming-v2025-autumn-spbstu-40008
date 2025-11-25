package internal

type ValuteXML struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type ResultValute struct {
	NumCode  int     `json:"num_code"  xml:"num_code"  yaml:"num_code"`
	CharCode string  `json:"char_code" xml:"char_code" yaml:"char_code"`
	Value    float64 `json:"value"     xml:"value"     yaml:"value"`
}
