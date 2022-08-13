package factory

import (
	"errors"

	"github.com/frizz925/covid19-update-bot/internal/country"
	"github.com/frizz925/covid19-update-bot/internal/fetcher"
	"github.com/frizz925/covid19-update-bot/internal/fetcher/factory"
	"github.com/frizz925/covid19-update-bot/internal/fetcher/jp/mhlw"
	"github.com/frizz925/covid19-update-bot/internal/scraper"
	jpScraper "github.com/frizz925/covid19-update-bot/internal/scraper/jp"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrNotImplemented = errors.New("not yet implemented")
	ErrInvalidFetcher = errors.New("invalid fetcher")
)

type ScraperFactory struct {
	factory.FetcherFactory
}

func NewScraperFactory(fixtureDir string) *ScraperFactory {
	return &ScraperFactory{
		FetcherFactory: factory.FetcherFactory{
			FixtureDir: fixtureDir,
		},
	}
}

func (f *ScraperFactory) Create(st scraper.Type, ft fetcher.Type, c country.Country, source string) (scraper.Scraper, error) {
	switch st {
	case scraper.Parsed:
		return f.ParsedScraper(ft, c, source)
	case scraper.Image:
		return f.ImageScraper(ft, c, source)
	}
	return nil, ErrNotFound
}

func (f *ScraperFactory) ParsedScraper(ft fetcher.Type, c country.Country, source string) (scraper.Scraper, error) {
	pf, err := f.FetcherFactory.ParsedFetcher(ft, c, source)
	if err != nil {
		return nil, err
	}
	return scraper.NewParsedScraper(pf), nil
}

func (f *ScraperFactory) ImageScraper(ft fetcher.Type, c country.Country, source string) (scraper.ImageScraper, error) {
	imgf, err := f.FetcherFactory.ImageFetcher(ft, c, source)
	if err != nil {
		return nil, err
	}
	switch c {
	case country.JP:
		if v, ok := imgf.(mhlw.Fetcher); ok {
			return jpScraper.NewMHLWScraper(v), nil
		}
		return nil, ErrInvalidFetcher
	case country.ID:
		return nil, ErrNotImplemented
	}
	return nil, ErrNotFound
}