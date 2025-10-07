package structures

type ValCurs struct {
	Valute []Valute `xml:"Valute"`
}
type Valute struct {
	NumCode  string `xml:"NumCode" json:"num_code"`
	CharCode string `xml:"CharCode" json:"char_code"`
	Value    string `xml:"Value" json:"value"`
}
