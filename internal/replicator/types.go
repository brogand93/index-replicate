package replicator

type Client struct {
	Index              string
	Percentage         float32
	Contribution       float64
	RoundShareQuantity bool
}

type ReplicatedIndex struct {
	Components []Component
}

type Component struct {
	name   string
	symbol string
	weight float32
	shares float32
}
