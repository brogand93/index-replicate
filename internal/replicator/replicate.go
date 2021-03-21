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

func (c *Client) Run(outputToCsv bool) error {
	indexClient := index.NewClient()
	index, err := indexClient.Get(c.Index)
	if err != nil {
		return err
	}

	percentageCovered := float32(0)
	totalContribution := float32(0)
	replicatedIndex := ReplicatedIndex{}

	for _, component := range index.Components {
		if (percentageCovered + component.Weight) < c.Percentage {
			weightInReplicatedIndex := c.componentWeight(component)
			sharesInRerplicatedIndex := c.sharesToBuy(component, weightInReplicatedIndex)
			valueInIndex := float32(component.Price) * sharesInRerplicatedIndex
			if c.RoundShareQuantity {
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

	err = c.write(replicatedIndex, outputToCsv)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// componentWeight calculates the weight of an individual component in
// the replciated index
func (c *Client) componentWeight(component index.Component) float32 {
	if component.Price > c.Contribution {
		return 0
	}

	return component.Weight * (100 / c.Percentage)
}

// sharesToBuy calculates the number of component shares to buy for an individual
// company in the replciated index
func (c *Client) sharesToBuy(component index.Component, weight float32) float32 {
	// Compute the portion of contribution that should be used on this component
	availableForComponent := float32(c.Contribution) / 100 * weight
	return availableForComponent / float32(component.Price)
}

func (c *Client) write(replicatedIndex ReplicatedIndex, outputToCsv bool) error {
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

	err := c.writeTable(data, replicatedIndex.Total)
	if err != nil {
		return err
	}

	if outputToCsv {
		return c.writeCSV(data)
	}

	return nil
}

func (c *Client) writeTable(data [][]string, total float32) error {
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

func (c *Client) writeCSV(data [][]string) error {
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
