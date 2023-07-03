package main

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	"github.com/jlgallego99/TropesToGo/service/crawler"
	"github.com/jlgallego99/TropesToGo/service/scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

var promptCrawlLimit = &survey.Input{
	Message: "How many Films would you like to extract?",
}

var promptLimit = &survey.Confirm{
	Message: "Do you want to scrape all Films?",
}

func main() {
	var crawlLimitInput string
	var crawlLimit int
	var unlimitedCrawling bool

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := survey.AskOne(promptLimit, &unlimitedCrawling)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	if unlimitedCrawling {
		crawlLimit = -1
		log.Info().Msg("Extracting all films in TvTropes...")
	} else {
		err = survey.AskOne(promptCrawlLimit, &crawlLimitInput, survey.WithValidator(numberValidator))
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		crawlLimit, _ = strconv.Atoi(crawlLimitInput)
		log.Info().Msg("Extracting " + crawlLimitInput + " films...")
	}

	start := time.Now()

	// Crawling TvTropes Pages
	serviceCrawler := crawler.NewCrawler()
	pages, err := serviceCrawler.CrawlWorkPages(crawlLimit)
	if err != nil && pages == nil {
		panic(err)
	}

	// Extracting data from TvTropes Pages and persisting them on a dataset file
	jsonRepository, _ := json_dataset.NewJSONRepository("dataset")
	serviceScraper, err := scraper.NewServiceScraper(scraper.ConfigMediaRepository(jsonRepository))
	if err != nil {
		panic(err)
	}

	errScraping := serviceScraper.ScrapeTvTropes(pages)
	if errScraping != nil {
		panic(errScraping)
	}

	log.Info().Msgf("Process finished in %s\n", time.Since(start))
	log.Info().Msg("TropesToGo finished successfully!")
}

func numberValidator(val interface{}) error {
	if _, err := strconv.Atoi(val.(string)); err != nil {
		return errors.New("the input must be a number")
	}

	return nil
}
