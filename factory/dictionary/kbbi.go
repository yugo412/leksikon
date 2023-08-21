package dictionary

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

type KBBI struct {
	URL      string
	Username string
	Password string

	Log *zap.SugaredLogger
}

func (*KBBI) Vendor() string {
	return "KBBI"
}

func (k *KBBI) login(c *colly.Collector) (err error) {
	parsedURL, err := url.Parse(strings.ReplaceAll(k.URL, "%s", ""))
	if err != nil {
		return err
	}

	host := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	loginURL, err := url.JoinPath(host, "Account", "Login")
	if err != nil {
		return
	}

	err = c.Post(loginURL, map[string]string{
		"Posel":     k.Username,
		"KataSandi": k.Password,
	})

	return
}

func (k *KBBI) Search(entry string) (words Word, err error) {
	entry = strings.ToLower(entry)

	c := colly.NewCollector()

	// remove duplicate spaces
	space := regexp.MustCompile(`\s+`)

	c.OnHTML("div.body-content", func(e *colly.HTMLElement) {
		if strings.Contains(e.Request.URL.String(), "Beranda/BatasSehari") {
			err = ErrMaxLimit

			return
		}

		// parse word & spell
		var root, spell, word, syllable string
		var informals []string
		e.DOM.Find("h2").Each(func(_ int, h2 *goquery.Selection) {
			syllable = h2.Find("span.syllable").First().Text()
			root = h2.Find("span.rootword > a").Text()

			// find multiple informal words
			h2.Find("small > b").Each(func(_ int, b *goquery.Selection) {
				text := strings.TrimSpace(strings.ReplaceAll(b.Text(), ",", ""))
				informals = append(informals, text)
			})

			w := h2.Find("span, sup, small").Remove().End()
			spell = strings.TrimSpace(w.Text())
			word = strings.ReplaceAll(spell, ".", "")
		})

		var definitions []WordDefinition
		var classes []string
		var examples []string
		var description string

		// multiple description for one word
		e.DOM.Find("ol > li").Each(func(i int, li *goquery.Selection) {
			li.Find("a.entrisButton").Remove().End()

			classes = []string{}
			class := space.ReplaceAllString(li.Find("font[color=red]").Text(), " ")

			examples = []string{}
			example := li.Find("font[color=grey]").Next().Text()
			if example != "" {
				examples = append(examples, strings.ReplaceAll(example, "~", "--"))
			}

			description = li.Find("font,i").Remove().End().Text()
			if description == "" {
				return
			}

			definitions = append(definitions, WordDefinition{
				Description: strings.TrimRight(description, ":"),
				Examples:    examples,
				Classes:     strings.Split(strings.TrimSpace(class), " "),
			})
		})

		e.DOM.Find("ul.adjusted-par > li").Each(func(i int, ul *goquery.Selection) {
			ul.Find("a").Remove().End()

			classes = []string{}
			class := space.ReplaceAllString(strings.TrimSpace(ul.Find("font").First().Text()), " ")
			if class != "" {
				classes = strings.Split(class, " ")
			}

			examples = []string{}
			example := ul.Find("font[color=grey]").Next().Text()
			if example != "" {
				examples = append(examples, strings.ReplaceAll(example, "~", "--"))
			}

			d := ul.Find("font").Remove().End()
			if strings.TrimSpace(strings.ReplaceAll(d.Text(), ";", "")) == "" {
				return
			}

			definitions = append(definitions, WordDefinition{
				Description: strings.ReplaceAll(d.Text(), ":", ""),
				Classes:     classes,
				Examples:    examples,
			})
		})

		words = Word{
			Root:        root,
			Word:        word,
			Spell:       spell,
			Syllable:    syllable,
			Informals:   informals,
			Source:      e.Request.URL.String(),
			Definitions: definitions,
		}
	})

	k.URL = strings.TrimSuffix(k.URL, "/")
	_ = c.Visit(fmt.Sprintf(k.URL, entry))

	return
}
