package schedule

import (
	"errors"
	"log/slog"
	"net/url"

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

func (s Scraper) GetFields(site string) map[string]string {
	f := make(map[string]string)

	s.c.OnHTML("ul.lista-grup a[href]", func(e *colly.HTMLElement) {
		l := e.Attr("href")
		n := e.Text
		u, err := url.Parse(l)
		if err != nil {
			return
		}
		id := u.Query().Get("ID")
		f[id] = n
	})

	if err := s.c.Visit(site); err != nil {
		slog.Warn("Failed to visit: %v", site)
	}
	return f
}

func (s Scraper) GetGroupsFromID(site string, id string) map[string]string {
	f := make(map[string]string)

	s.c.OnHTML("table a[href]", func(e *colly.HTMLElement) {
		l := e.Attr("href")
		n := e.Text
		u, err := url.Parse(l)
		if err != nil {
			return
		}
		id := u.Query().Get("ID")
		f[id] = n
	})

	u, _ := url.Parse(site)
	q := u.Query()
	q.Set("ID", id)
	u.RawQuery = q.Encode()
	if err := s.c.Visit(u.String()); err != nil {
		slog.Warn("Failed to visit: %v", site)
	}
	return f
}

func (s Scraper) GetScheduleForId(site string, id string) ([]Entry, error) {
	var entries []Entry
	var errs []error

	s.c.OnHTML("#table_details tr:has(td)", func(e *colly.HTMLElement) {
		date := e.ChildText("td:nth-child(1)")
		_ = e.ChildText("td:nth-child(2)")
		group := e.ChildText("td:nth-child(3)")
		start := e.ChildText("td:nth-child(4)")
		end := e.ChildText("td:nth-child(5)")
		subject := e.ChildText("td:nth-child(6)")
		ClassType := e.ChildText("td:nth-child(7)")
		teacher := e.ChildText("td:nth-child(8)")
		classroom := e.ChildText("td:nth-child(9)")
		ent, err := FromScraper(RawEntry{
			Group:     group,
			Start:     start,
			End:       end,
			Date:      date,
			Subject:   subject,
			ClassType: ClassType,
			Teacher:   teacher,
			Classroom: classroom,
		})
		if err != nil {
			errs = append(errs, err)
		}
		entries = append(entries, ent)
	})

	u, _ := url.Parse(site)
	q := u.Query()
	q.Set("ID", id)
	u.RawQuery = q.Encode()
	if err := s.c.Visit(u.String()); err != nil {
		slog.Warn("Failed to visit", site)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return entries, nil
}
