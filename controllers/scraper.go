package controllers

import (
	"github.com/gocolly/colly"
)

// SoundCloudScraper is a scraper for SoundCloud.
type Scraper struct {
	collector *colly.Collector
}

func NewScraper() *Scraper {
	return &Scraper{
		collector: colly.NewCollector(),
	}
}

func (s *Scraper) Scrape(url string) error {
	s.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		s.collector.Visit(link)
	})

	return s.collector.Visit(url)
}
