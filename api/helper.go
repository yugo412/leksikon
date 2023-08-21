package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
)

var serverTime = time.Now().Local().Format(DefaultDateTimeFormat)

type RootResponse struct {
	Status     string      `json:"message"`
	ServerTime string      `json:"server_time"`
	Data       interface{} `json:"data"`
}

const (
	DefaultDateTimeFormat = time.RFC3339
)

func OK(w http.ResponseWriter, content interface{}) (err error) {
	var res RootResponse
	res.Status = "OK"
	res.ServerTime = serverTime
	res.Data = content

	body, err := json.Marshal(res)
	if err != nil {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)

	return
}

func OKWithStatus(w http.ResponseWriter, status int, content interface{}) (err error) {
	if status == http.StatusNoContent {
		_, err = w.Write([]byte{})

		return
	}

	var res RootResponse
	res.Status = "OK"
	res.ServerTime = serverTime
	res.Data = content

	body, err := json.Marshal(res)
	if err != nil {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(body)

	return
}

func Error(w http.ResponseWriter, status int, content interface{}) (err error) {
	var res RootResponse
	res.ServerTime = serverTime

	if reflect.TypeOf(content).Kind() == reflect.String {
		res.Status = content.(string)
	} else {
		res.Status = fmt.Sprintf("%v", content)

		if reflect.TypeOf(validator.ValidationErrors{}) == reflect.TypeOf(content) {
			validations := make(map[string][]string)
			for _, v := range content.(validator.ValidationErrors) {
				validations[v.Field()] = append(validations[v.Field()], v.Tag())

			}

			res.Status = "failed to validate request"
			res.Data = validations
		}
	}

	body, err := json.Marshal(res)
	if err != nil {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(body)

	return
}
