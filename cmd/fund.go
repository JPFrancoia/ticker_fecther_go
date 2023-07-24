// Fetches information for index funds.
// Uses the Yahoo API.

package cmd

import (
	"fmt"
	"local/ticker_fetcher/utils"
	"sort"
	"sync"

	"local/ticker_fetcher/yahoo"

	"github.com/spf13/cobra"
)

var fundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Fetch performances for ETFs",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		funds := utils.ParseCommaString(args[0])

		c := make(chan yahoo.YahooInfo)
		var wg sync.WaitGroup

		for _, fund := range funds {
			wg.Add(1)
			go yahoo.FetchInfoFromYahoo(fund, c, &wg)
		}

		// Necessary to avoid deadlock
		// https://stackoverflow.com/a/70877210/1585507
		go func() {
			wg.Wait()
			close(c)
		}()

		var results []yahoo.YahooInfo

		for data := range c {
			results = append(results, data)
		}

		// Ascending sort
		sort.Slice(results, func(i, j int) bool {
			return results[i].Symbol < results[j].Symbol
		})

		for _, data := range results {
			fmt.Printf(
				"${alignc}%s: %g %s (%.2f %%)\n",
				data.Symbol,
				data.Price,
				data.Currency,
				data.Diff(),
			)
		}
	},
}
