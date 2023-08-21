package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/yugo412/leksikon/api/handler"
	"github.com/yugo412/leksikon/model"
	"github.com/yugo412/leksikon/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func config() error {
	viper.SetConfigFile(".env")

	return viper.ReadInConfig()
}

func logger() (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error

	if viper.GetString("ENV") == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	defer logger.Sync()

	return logger.Sugar(), err

}

func database() (*gorm.DB, error) {
	config := mysql.Config{
		User:                 viper.GetString("DB_USER"),
		Passwd:               viper.GetString("DB_PASSWORD"),
		Net:                  viper.GetString("DB_PROTOCOL"),
		Addr:                 fmt.Sprintf("%s:%s", viper.GetString("DB_HOST"), viper.GetString("DB_PORT")),
		DBName:               viper.GetString("DB_NAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := gorm.Open(gormMySQL.Open(config.FormatDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if viper.GetBool("DB_DEBUG") {
		db = db.Debug()
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Word{},
		&model.Definition{},
		&model.Example{},
		&model.Language{},
		&model.Slang{},
		&model.Class{},
		&model.Counter{},
		&model.Source{},
	)
}

var (
	orm *gorm.DB
	log *zap.SugaredLogger
)

func init() {
	var err error

	err = config()
	if err != nil {
		fmt.Println("file .env does not exist")
		os.Exit(2)
	}

	log, err = logger()
	if err != nil {
		fmt.Println("cannot init logger: ", err)
		os.Exit(1)
	}

	orm, err = database()
	if err != nil {
		fmt.Println("cannot connect to database: ", err)
		os.Exit(1)
	}

	if viper.GetString("ENV") != "production" {
		err = migrate(orm)
		if err != nil {
			fmt.Println("cannot migrate database: ", err)
			os.Exit(1)
		}
	}
}

func main() {
	var err error

	r := chi.NewRouter()
	if viper.GetString("ENV") != "production" {
		r.Use(middleware.Logger)
	}

	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.AllowContentType("application/json", "text/xml"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var res struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}

		res.Name = viper.GetString("NAME")
		res.Version = viper.GetString("VERSION")

		body, _ := json.Marshal(res)
		_, _ = w.Write(body)
	})

	lang := handler.NewLanguage(service.Language{
		ORM: orm,
		Log: log,

		Lang: &model.LanguageModel{ORM: orm},
	})
	r.Get("/language", lang.GetAllLang)

	dict := handler.NewDictionary(service.Dictionary{
		ORM: orm,
		Log: log,

		Word:   &model.WordModel{ORM: orm},
		Lang:   &model.LanguageModel{ORM: orm},
		Source: &model.SourceModel{DB: orm},
	})
	r.Get("/dictionary", dict.Index)
	r.Get("/dictionary/entry", dict.Search)
	r.Get("/dictionary/live", dict.LiveSearch)

	err = http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("PORT")), r)
	if err != nil {
		panic(err)
	}
}
