package valute

import (

)

type StructOfXMLandJSON struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
}

type ValuteCurs struct {

}

type Valute struct {
	XMLName xml.Name               `xml:"ValCurs"`
	Valutes []StructOfXMLandJSON   `xml:"Valute"`
}

type Convertion struct {

}

