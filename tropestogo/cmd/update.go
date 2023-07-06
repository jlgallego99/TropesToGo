package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jlgallego99/TropesToGo/media"
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	"github.com/jlgallego99/TropesToGo/service/crawler"
	"github.com/jlgallego99/TropesToGo/service/scraper"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var (
	updateDatasetName, updateFileFormat string

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Updates an already-extracted dataset with new updated data, if there's any on TvTropes",
		Long: `The update command updates the local dataset file by providing its name with the -d flag.
By default, searches for a file with the name "dataset.json".`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Msg("Launching TropesToGo Updater")

			_, errFileExists := os.Stat(datasetPath + "/" + updateDatasetName)
			if errFileExists != nil {
				log.Error().Err(errFileExists).Msg("Couldn't retrieve the dataset file " + updateDatasetName + " on " + datasetPath)
				return errFileExists
			}

			scrapeUpdates()

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringVarP(&updateDatasetName, "dataset", "d", "dataset.json", "must specify a name for the dataset to update with the extension (-d <datasetfile>)")
}

func scrapeUpdates() {
	start := time.Now()

	// Call scraper to extract the persisted changedPages on the dataset
	var repository media.RepositoryMedia
	updateFileFormat = strings.ReplaceAll(filepath.Ext(updateDatasetName), ".", "")
	datasetBaseName := strings.TrimSuffix(updateDatasetName, filepath.Ext(updateDatasetName))
	if strings.EqualFold(updateFileFormat, CSV) {
		repository, _ = csv_dataset.NewCSVRepository(datasetBaseName)
	} else if strings.EqualFold(updateFileFormat, JSON) {
		repository, _ = json_dataset.NewJSONRepository(datasetBaseName)
	}

	serviceScraper, err := scraper.NewServiceScraper(scraper.ConfigMediaRepository(repository))
	if err != nil {
		log.Error().Err(err).Msg("Error creating TropesToGo scraper")
		return
	}

	pagesToBeUpdated, errScrapedPages := serviceScraper.GetScrapedPages()
	if errScrapedPages != nil {
		log.Error().Err(errScrapedPages).Msg("Scraping changedPages error")
		return
	}

	// Crawling Pages with updates
	serviceCrawler := crawler.NewCrawler()
	changedPages, err := serviceCrawler.CrawlChanges(pagesToBeUpdated)
	if err != nil && changedPages == nil {
		log.Error().Err(err).Msg("Error in TropesToGo crawling changes")
		return
	}

	// Updating changedPages
	if len(changedPages.Pages) > 0 {
		serviceScraper.UpdateDataset(changedPages)

		log.Info().Msg(strconv.Itoa(len(changedPages.Pages)) + " works have been updated in the dataset " + updateDatasetName)
		log.Info().Msg("The updated TvTropes dataset is available on: " + datasetPath + "service/" + updateDatasetName)
	} else {
		log.Info().Msg("The dataset " + updateDatasetName + " is already up to date!")
	}

	log.Info().Msgf("Process finished in %s\n", time.Since(start))
	log.Info().Msg("TropesToGo finished successfully!")
}
