package main

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/jlgallego99/TropesToGo/media/json_dataset"
	"github.com/jlgallego99/TropesToGo/service/crawler"
	"github.com/jlgallego99/TropesToGo/service/scraper"
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

	err := survey.AskOne(promptLimit, &unlimitedCrawling)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if unlimitedCrawling {
		crawlLimit = -1
		fmt.Println("Extracting all films in TvTropes...")
	} else {
		err = survey.AskOne(promptCrawlLimit, &crawlLimitInput, survey.WithValidator(numberValidator))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		crawlLimit, _ = strconv.Atoi(crawlLimitInput)
		fmt.Println("Extracting", crawlLimit, "films...")
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

	fmt.Printf("Process finished in %s\n", time.Since(start))
	fmt.Println("TropesToGo finished successfully!")
}

func numberValidator(val interface{}) error {
	if _, err := strconv.Atoi(val.(string)); err != nil {
		return errors.New("the input must be a number")
	}

	return nil
}
