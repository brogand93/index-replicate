package replicator

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/brogand93/index-replicate/pkg/index"
	"github.com/olekukonko/tablewriter"
)

const (
	maxColWidth = 50
	outputCSV   = "index_info.csv"
)

func (client *Client) Run(outputToCsv bool) error {
	index, err := index.Get(client.Index)
	if err != nil {
		return err
	}

	percentageCovered := float32(0)
	totalContribution := float32(0)
	replicatedIndex := ReplicatedIndex{}

	for _, component := range index.Components {
		if (percentageCovered + component.Weight) < client.Percentage {
			weightInReplicatedIndex := client.componentWeight(component)
			sharesInRerplicatedIndex := client.sharesToBuy(component, weightInReplicatedIndex)
			valueInIndex := float32(component.Price) * sharesInRerplicatedIndex
			if client.RoundShareQuantity {
				sharesInRerplicatedIndex = float32(
					math.Round(float64(sharesInRerplicatedIndex)),
				)
			}
			replicatedIndex.Components = append(replicatedIndex.Components, Component{
				Name:   component.Name,
				Symbol: component.Symbol,
				Weight: weightInReplicatedIndex,
				Shares: sharesInRerplicatedIndex,
				Value:  valueInIndex,
			})
			totalContribution = totalContribution + valueInIndex
		}

		percentageCovered = percentageCovered + component.Weight
	}

	replicatedIndex.Total = totalContribution

	err = client.write(replicatedIndex, outputToCsv)
	if err != nil {
		log.Println(err.Error())
		return err
	}

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

func (client *Client) write(replicatedIndex ReplicatedIndex, outputToCsv bool) error {
	data := [][]string{}

	for _, component := range replicatedIndex.Components {
		data = append(data, []string{
			component.Name,
			component.Symbol,
			fmt.Sprintf("%.2f", component.Shares),
			fmt.Sprintf("%.2f %%", component.Weight),
			fmt.Sprintf("%.2f $", component.Value),
		})
	}

	err := client.writeTable(data, replicatedIndex.Total)
	if err != nil {
		return err
	}

	if outputToCsv {
		return client.writeCSV(data)
	}

	return nil
}

func (client *Client) writeTable(data [][]string, total float32) error {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{
		"Name", "Symbol", "Quantity", "Weight", "Value",
	})
	table.SetFooter([]string{
		"", "", "", "Total", fmt.Sprintf("%.2f $", total),
	})
	table.SetColWidth(maxColWidth)
	table.AppendBulk(data)
	table.Render()
	return nil

}

func (client *Client) writeCSV(data [][]string) error {
	file, err := os.Create(outputCSV)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}
