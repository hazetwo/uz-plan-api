package schedule

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
}

func NewScraper() *Scraper {
	return &Scraper{}
}

func getDocument(site string) (*goquery.Document, error) {
	res, err := http.Get(site)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func getURLWithID(site string, id string) (string, error) {
	u, err := url.Parse(site)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("ID", id)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (s Scraper) GetFields(site string) (map[string]string, error) {
	f := make(map[string]string)

	doc, err := getDocument(site)
	if err != nil {
		return nil, err
	}

	doc.Find("ul.lista-grup a[href]").Each(func(i int, s *goquery.Selection) {
		l, _ := s.Attr("href")
		n := s.Text()
		u, err := url.Parse(l)
		if err != nil {
			return
		}
		id := u.Query().Get("ID")
		f[id] = n
	})

	return f, nil
}

func (s Scraper) GetGroupsFromID(site string, id string) (map[string]string, error) {
	g := make(map[string]string)

	u, err := getURLWithID(site, id)
	if err != nil {
		return nil, err
	}

	doc, err := getDocument(u)
	if err != nil {
		return nil, err
	}

	doc.Find("table a[href]").Each(func(i int, s *goquery.Selection) {
		l, ok := s.Attr("href")
		if !ok {
			return
		}
		n := s.Text()
		u, err := url.Parse(l)
		if err != nil {
			return
		}
		id := u.Query().Get("ID")
		g[id] = n
	})

	return g, nil
}

func (s Scraper) GetScheduleForID(site string, id string) ([]Entry, error) {
	var entries []Entry
	var errs []error

	u, err := getURLWithID(site, id)
	if err != nil {
		return nil, err
	}

	doc, err := getDocument(u)
	if err != nil {
		return nil, err
	}

	doc.Find("#table_details tr:has(td)").Each(func(i int, s *goquery.Selection) {
		date := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
		// WE SKIP THE 2nd child
		group := strings.TrimSpace(s.Find("td:nth-child(3)").Text())
		start := strings.TrimSpace(s.Find("td:nth-child(4)").Text())
		end := strings.TrimSpace(s.Find("td:nth-child(5)").Text())
		subject := strings.TrimSpace(s.Find("td:nth-child(6)").Text())
		classType := strings.TrimSpace(s.Find("td:nth-child(7)").Text())
		teacher := strings.TrimSpace(s.Find("td:nth-child(8)").Text())
		classroom := strings.TrimSpace(s.Find("td:nth-child(9)").Text())
		e, err := FromScraper(RawEntry{
			Group:     group,
			Start:     start,
			End:       end,
			Date:      date,
			Subject:   subject,
			ClassType: classType,
			Teacher:   teacher,
			Classroom: classroom,
		})
		if err != nil {
			errs = append(errs, err)
			slog.Error("Errors caught when scraping schedule", "err", err)
		}
		entries = append(entries, e)
	})

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return entries, nil

}
