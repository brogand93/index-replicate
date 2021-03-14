package replicator

import (
	"fmt"
	"math"

	"github.com/brogand93/index-replicate/pkg/index"
)

func (client *Client) Run() error {
	index, err := index.Get(client.Index)
	if err != nil {
		return err
	}

	percentageCovered := float32(0)
	replicatedIndex := ReplicatedIndex{}

	for _, component := range index.Components {
		if (percentageCovered + component.Weight) < client.Percentage {
			weightInReplicatedIndex := client.componentWeight(component)
			sharesInRerplicatedIndex := client.sharesToBuy(component, weightInReplicatedIndex)
			if client.RoundShareQuantity {
				sharesInRerplicatedIndex = float32(
					math.Round(float64(sharesInRerplicatedIndex)),
				)
			}
			replicatedIndex.Components = append(replicatedIndex.Components, Component{
				name:   component.Name,
				symbol: component.Symbol,
				weight: weightInReplicatedIndex,
				shares: sharesInRerplicatedIndex,
			})
		}
		percentageCovered = percentageCovered + component.Weight
	}

	client.write(replicatedIndex)

	return nil
}

// componentWeight calculates the weight of an individual component in
// the replciated index
func (client *Client) componentWeight(component index.Component) float32 {
	if component.Price > client.Contribution {
		return 0
	}

	return component.Weight * (100 / client.Percentage)
}

// sharesToBuy calculates the number of component shares to buy for an individual
// company in the replciated index
func (client *Client) sharesToBuy(component index.Component, weight float32) float32 {
	// Compute the portion of contribution that should be used on this component
	availableForComponent := float32(client.Contribution) / 100 * weight
	return availableForComponent / float32(component.Price)
}

func (client *Client) write(replicatedIndex ReplicatedIndex) {
	for _, component := range replicatedIndex.Components {
		fmt.Printf("%s : %v\n", component.symbol, component.shares)
	}
}
