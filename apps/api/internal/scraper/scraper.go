package scraper

import (
	"net/url"
	"slices"

	"github.com/gocolly/colly"
)

type Scraper struct {
	c *colly.Collector
}

func New(domain string) Scraper {
	return Scraper{
		c: colly.NewCollector(
			colly.AllowedDomains(domain)),
	}
}

func (s Scraper) GetIdsOfFields(site string) []string {
	var ids []string

	s.c.OnHTML("ul.lista-grup a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		u, err := url.Parse(link)
		if err != nil {
			return
		}
		id := u.Query().Get("ID")
		ids = append(ids, id)
	})

	s.c.Visit(site)
	return ids
}

func (s Scraper) GetIdsOfGroups(site string, fields []string, supportedIds []string) []string {
	var ids []string

	s.c.OnHTML("table a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		u, err := url.Parse(link)
		if err != nil {
			return
		}
		id := u.Query().Get("ID")
		ids = append(ids, id)
	})

	for _, f := range fields {
		if !slices.Contains(supportedIds, f) {
			continue
		}
		u, _ := url.Parse(site)
		q := u.Query()
		q.Set("ID", f)
		u.RawQuery = q.Encode()
		s.c.Visit(u.String())
	}

	return ids
}
