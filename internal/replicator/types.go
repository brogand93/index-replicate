package replicator

type Client struct {
	RoundShareQuantity bool
	Percentage         float32
	Contribution       float64
	Index              string
}

type ReplicatedIndex struct {
	Total      float32
	Components []Component
}

type Component struct {
	Name   string
	Symbol string
	Weight float32
	Shares float32
	Value  float32
}
