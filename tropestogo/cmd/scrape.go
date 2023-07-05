package cmd

import (
	"fmt"
	"github.com/jlgallego99/TropesToGo/media"
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	"github.com/jlgallego99/TropesToGo/service/crawler"
	"github.com/jlgallego99/TropesToGo/service/scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CSV  string = "CSV"
	JSON        = "JSON"
)

// scrapeCmd represents the scrape command
var (
	datasetPath, _          = os.Getwd()
	datasetName, dataFormat string
	crawlLimit              int
	crawlAll                bool

	scrapeCmd = &cobra.Command{
		Use:   "scrape",
		Short: "Scrapes works of any media type with its tropes and generates a dataset",
		Long: `The scrape command is the main TropesToGo command for scraping works 
of any media type with its tropes from TvTropes.
Generates a dataset of the specified format when done.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logFile, _ := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
			multiWriter := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr}, logFile)
			log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
			log.Info().Msg("TropesToGo: A scraper for TvTropes")

			if !strings.EqualFold(dataFormat, CSV) && !strings.EqualFold(dataFormat, JSON) {
				return fmt.Errorf("unknown data format: %s", dataFormat)
			}

			if crawlAll {
				log.Info().Msg("Extracting all films in TvTropes...")
				crawlLimit = -1
			} else {
				log.Info().Msg("Extracting " + strconv.Itoa(crawlLimit) + " films...")
			}

			scrape()

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.PersistentFlags().StringVarP(&datasetName, "output", "o", "dataset", "specify a name for the dataset (-o <datasetname>)")
	scrapeCmd.PersistentFlags().StringVarP(&dataFormat, "format", "f", "json", "specify a format for the dataset (-f json, -f csv)")
	scrapeCmd.PersistentFlags().IntVarP(&crawlLimit, "limit", "l", 1, "limit the number of extracted works (-l <number>)")
	scrapeCmd.PersistentFlags().BoolVarP(&crawlAll, "all", "a", false, "if set, it extracts all works, on the contrary it will extract the number specified with the -l flag")
}

func scrape() {
	start := time.Now()

	// Crawling TvTropes Pages
	serviceCrawler := crawler.NewCrawler()
	pages, err := serviceCrawler.CrawlWorkPages(crawlLimit)
	if err != nil && pages == nil {
		log.Error().Err(err).Msg("Error creating TropesToGo crawler")
		return
	}

	// Extracting data from TvTropes Pages and persisting them on a dataset file
	var repository media.RepositoryMedia
	if strings.EqualFold(dataFormat, CSV) {
		repository, _ = csv_dataset.NewCSVRepository(datasetName)
		datasetName += "." + strings.ToLower(CSV)
	} else if strings.EqualFold(dataFormat, JSON) {
		repository, _ = json_dataset.NewJSONRepository(datasetName)
		datasetName += "." + strings.ToLower(JSON)
	}

	serviceScraper, err := scraper.NewServiceScraper(scraper.ConfigMediaRepository(repository))
	if err != nil {
		log.Error().Err(err).Msg("Error creating TropesToGo scraper")
		return
	}

	errScraping := serviceScraper.ScrapeTvTropes(pages)
	if errScraping != nil {
		log.Error().Err(errScraping).Msg("Scraping error")
		return
	}

	log.Info().Msgf("Process finished in %s\n", time.Since(start))
	log.Info().Msg("TropesToGo finished successfully!")
	log.Info().Msg("The generated TvTropes dataset is available on: " + datasetPath + "service/" + datasetName)
}
