package structures

type File struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

type ValCursXML struct {
	Valute []ValuteXML `xml:"Valute"`
}
type ValuteXML struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type ValCursJSON struct {
	Valute []ValuteJSON
}
type ValuteJSON struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}
