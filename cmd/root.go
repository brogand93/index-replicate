package cmd

import (
	"github.com/brogand93/index-replicate/internal/replicator"
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "index-replicate",
	Short: "Index Replicate is a lightewight client to pull stock index",
	Long:  `Index Replicate allows you to replicate an index`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	replicator := replicator.Client{
		Index:              indexToReplicate,
		Contribution:       contribution,
		Percentage:         float32(percentage),
		RoundShareQuantity: roundShareQuantity,
	}

	replicator.Run()
}

var indexToReplicate string
var contribution float64
var percentage int
var roundShareQuantity bool

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&indexToReplicate, "index", "i", "sp500", "index to replicate [sp500, dowjones, nasdaq100]")
	rootCmd.MarkFlagRequired("index")
	rootCmd.Flags().Float64VarP(&contribution, "contribution", "c", 0, "amount you wish to invest in the index")
	rootCmd.MarkFlagRequired("contribution")
	rootCmd.Flags().IntVarP(&percentage, "percentage", "p", 100, "percentage of the index you wish to replicate")
	rootCmd.Flags().BoolVarP(&roundShareQuantity, "round", "r", false, "round share buy quantity to the nearest whole share")
}
