package library

import (
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type LocalLanguage struct {
	Name   string
	Region string
}

func FetchLocalLanguage() (languages []LocalLanguage, err error) {
	c := colly.NewCollector()

	c.OnHTML("table.isi > tbody > tr", func(el *colly.HTMLElement) {
		name := el.ChildText("td:nth-child(2)")
		region := el.ChildText("td:nth-child(3)")

		// sanitize the name
		name = strings.ToValidUTF8(name, "")
		region = strings.ToValidUTF8(region, "")

		// remove double space for name
		space, _ := regexp.Compile(`\s+`)
		name = space.ReplaceAllString(name, " ")

		// remove double space for region
		region = space.ReplaceAllString(region, " ")

		languages = append(languages, LocalLanguage{
			Name:   name,
			Region: region,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	err = c.Visit("https://petabahasa.kemdikbud.go.id/databahasa.php")
	if err != nil {
		return
	}

	return
}
