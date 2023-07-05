package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tropestogo",
	Short: "TropesToGo - a scraper for TvTropes",
	Long: `TropesToGo is a scraper that can extract all works of any media type with its associated tropes from TvTropes.
It generates a dataset with all the scraped data. 
Examples of use:

- tropestogo scrape -o mydataset -f csv -l 10
this will extract 10 works with its tropes from TvTropes, and store them on a mydataset.csv file`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	logFile, _ := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	multiWriter := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr}, logFile)
	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
	log.Info().Msg("TropesToGo: A scraper for TvTropes")

	err := rootCmd.Execute()
	if err != nil {
		log.Error().Err(err).Msg("There was a problem on the CLI program")
	}
}

func init() {
}
