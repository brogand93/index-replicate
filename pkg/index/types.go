package index

import "net/http"

type Client struct {
	Http HTTPClient
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Index represents a stcok index
// an index is made up of multiple companies (represented as components of the index)
type Index struct {
	Components []Component `json:"companyList"`
}

// Component represents a position in an index
type Component struct {
	Name   string  `json:"name"`
	Price  float64 `json:"lastPrice"`
	Rank   float64 `json:"rank"`
	Symbol string  `json:"symbol"`
	Weight float32 `json:"weight"`
	Value  float32 `json:"value"`
}

type source struct {
	url string
}

type NotFound struct {
	Message string
}

type DataSourceChanged struct {
	Message string
}

type DataSourceUnavailable struct {
	Message string
}
