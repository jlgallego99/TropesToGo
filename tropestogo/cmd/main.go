package main

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/jlgallego99/TropesToGo/media"
	"github.com/jlgallego99/TropesToGo/media/csv_dataset"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	"github.com/jlgallego99/TropesToGo/service/crawler"
	"github.com/jlgallego99/TropesToGo/service/scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CSV  string = "CSV"
	JSON        = "JSON"
)

var datasetPath, _ = os.Getwd()

var promptDataFormat = &survey.Select{
	Message: "Select the data format for the generated dataset",
	Options: []string{"CSV", "JSON"},
}

var promptDatasetName = &survey.Input{
	Message: "What will be the name of the generated TvTropes dataset?",
}

var promptMediaType = &survey.Select{
	Message: "Select the media you wish to extract the data from",
	Options: media.GetAllMediaTypes(),
}

var promptLimit = &survey.Confirm{
	Message: "Do you want to scrape all works on the selected Media?",
}

var promptCrawlLimit = &survey.Input{
	Message: "How many works would you like to extract?",
}

func main() {
	var crawlLimitInput, datasetName, dataFormat, mediaTypeString string
	var mediaType media.MediaType
	var crawlLimit int
	var unlimitedCrawling bool

	logFile, _ := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	multiWriter := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr}, logFile)
	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
	log.Info().Msg("TropesToGo: A scraper for TvTropes")

	err := survey.AskOne(promptDataFormat, &dataFormat)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	err = survey.AskOne(promptDatasetName, &datasetName, survey.WithValidator(survey.Required))
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	err = survey.AskOne(promptMediaType, &mediaTypeString)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	mediaType, _ = media.ToMediaType(mediaTypeString)

	err = survey.AskOne(promptLimit, &unlimitedCrawling, survey.WithValidator(survey.Required))
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	if unlimitedCrawling {
		crawlLimit = -1
		log.Info().Msg("Extracting all works of type " + mediaTypeString + " in TvTropes...")
	} else {
		err = survey.AskOne(promptCrawlLimit, &crawlLimitInput, survey.WithValidator(numberValidator))
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		crawlLimit, _ = strconv.Atoi(crawlLimitInput)
		log.Info().Msg("Extracting " + crawlLimitInput + " works of type " + mediaTypeString + "...")
	}

	start := time.Now()

	// Crawling TvTropes Pages
	serviceCrawler := crawler.NewCrawler()
	pages, err := serviceCrawler.CrawlWorkPages(crawlLimit, mediaType)
	if err != nil && pages == nil {
		panic(err)
	}

	// Extracting data from TvTropes Pages and persisting them on a dataset file
	var repository media.RepositoryMedia
	switch {
	case dataFormat == CSV:
		repository, _ = csv_dataset.NewCSVRepository(datasetName)
		datasetName += "." + strings.ToLower(CSV)
	case dataFormat == JSON:
		repository, _ = json_dataset.NewJSONRepository(datasetName)
		datasetName += "." + strings.ToLower(JSON)
	default:
		repository, _ = json_dataset.NewJSONRepository(datasetName)
		datasetName += "." + strings.ToLower(JSON)
	}

	serviceScraper, err := scraper.NewServiceScraper(scraper.ConfigMediaRepository(repository))
	if err != nil {
		panic(err)
	}

	errScraping := serviceScraper.ScrapeTvTropes(pages)
	if errScraping != nil {
		panic(errScraping)
	}

	log.Info().Msgf("Process finished in %s\n", time.Since(start))
	log.Info().Msg("TropesToGo finished successfully!")
	log.Info().Msg("The generated TvTropes dataset is available on: " + datasetPath + "service/" + datasetName)
}

func numberValidator(val interface{}) error {
	var limitNumber int
	var errNumber error

	if limitNumber, errNumber = strconv.Atoi(val.(string)); errNumber != nil {
		return errors.New("the input must be a number")
	}

	if limitNumber <= 0 {
		return errors.New("there should be at least one work to scrape")
	}

	return nil
}
