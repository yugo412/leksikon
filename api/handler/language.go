package handler

import (
	"net/http"

	"github.com/yugo412/leksikon/api"
	"github.com/yugo412/leksikon/api/request"
	"github.com/yugo412/leksikon/service"

	"go.uber.org/zap"
)

type Language struct {
	Service service.Language
	Log     *zap.SugaredLogger
}

func NewLanguage(service service.Language) *Language {
	return &Language{
		Service: service,
		Log:     service.Log,
	}
}

func (l *Language) GetAllLang(w http.ResponseWriter, r *http.Request) {
	var req request.LanguageRequest
	req.Name = r.URL.Query().Get("name")

	languages, err := l.Service.GetAllLanguage(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = api.OK(w, languages)
	if err != nil {
		l.Log.Errorln(err)
	}
}
