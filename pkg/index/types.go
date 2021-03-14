package index

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

type IndexNotFound struct {
	message string
}

type DataSourceChanged struct {
	message string
}

type DataSourceUnavailable struct {
	message string
}
