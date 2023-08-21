package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/yugo412/leksikon/api"
	"github.com/yugo412/leksikon/factory/dictionary"
	"github.com/yugo412/leksikon/service"

	"go.uber.org/zap"
)

type Dictionary struct {
	Service service.Dictionary
	Log     *zap.SugaredLogger
}

func NewDictionary(service service.Dictionary) *Dictionary {
	return &Dictionary{
		Service: service,
		Log:     service.Log,
	}
}

func (d *Dictionary) LiveSearch(w http.ResponseWriter, r *http.Request) {
	entry := r.URL.Query().Get("entry")
	word, err := dictionary.New("https://kbbi.kemdikbud.go.id/entri/%s").Search(entry)
	if err != nil {
		_ = api.Error(w, http.StatusInternalServerError, err)

		return
	}

	if word.Word == "" {
		_ = api.Error(w, http.StatusNotFound, errors.New("word not found"))

		return
	}

	if len(word.Definitions) <= 0 {
		d.Log.Warnw("failed to fetch word definitions", "entry", entry)
	}

	_ = api.OK(w, word)

}

func (d *Dictionary) Index(w http.ResponseWriter, r *http.Request) {
	page := func() int {
		p := r.URL.Query().Get("page")
		if p == "" {
			return 1
		}

		i, _ := strconv.Atoi(p)

		return i
	}()

	entries, err := d.Service.Index(page)
	if err != nil {
		_ = api.Error(w, http.StatusInternalServerError, err)

		return
	}

	_ = api.OK(w, entries)
}

func (d *Dictionary) Search(w http.ResponseWriter, r *http.Request) {
	word, err := d.Service.Search(r.URL.Query().Get("entry"))
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrEntryNotFound) {
			code = http.StatusNotFound
		}
		_ = api.Error(w, code, err)

		return
	}

	_ = api.OK(w, word)
}
