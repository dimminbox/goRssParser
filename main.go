package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goRssParser/core"
	db2 "goRssParser/db"
	"goRssParser/rss"
	"log"
	"os"
	"runtime/debug"
	"time"
)

func connectDB(config db2.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.Database.Username,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Db,
		))
	db.SetMaxOpenConns(40)
	db.SetMaxIdleConns(5)
	if err != nil {
		panic(err)
	}
	return db, err
}
func main() {

	file, err := os.OpenFile("logger.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	errT := file.Truncate(0)
	if errT != nil {
		log.Fatal(errT)
	}

	defer file.Close()
	log.SetOutput(file)

	defer func() {
		if r := recover(); r != nil {
			log.Print("stacktrace from panic: \n" + string(debug.Stack()))
			log.Println("Recovered in f ", r)
		}
	}()

	config := db2.LoadConfiguration("config.json")

	timeout := 5
	if config.Sleep > 0 {
		timeout = config.Sleep
	}
	db, err := connectDB(config)

	if err != nil {
		panic(err)
	}

	for true {

		if config.Debug {
			log.Printf("Open sql connection \n")
		}

		for db.Ping() != nil {
			db, _ = connectDB(config)
			time.Sleep(time.Duration(timeout) * time.Second)
		}

		if config.Debug {
			log.Printf("Begin search keywords \n")
		}

		keywordList, errK := core.GetKeyword(db)
		if errK != nil {
			panic(errK)
		}

		if config.Debug {
			log.Printf("End search keywords \n")
		}

		if config.Debug {
			log.Printf("Begin search keywords \n")
		}

		keywordExList, errKe := core.GetKeywordExc(db)
		if errKe != nil {
			panic(errKe)
		}

		if config.Debug {
			log.Printf("End search keywords \n")
		}

		if config.Debug {
			log.Printf("Begin search website \n")
		}

		website, err2 := core.GetUrl(db, 15*60)
		if err2 != nil {
			panic(err2)
		}

		if config.Debug {
			log.Printf("End search website \n")
		}

		if website.Id > 0 {
			db.SetConnMaxLifetime(time.Second * time.Duration(website.Timeout))
			list := rss.ParseRss(website.Url, keywordList, keywordExList)

			for _, item := range list {
				if config.Debug {
					log.Printf("Begin save to DB %s \n", website.Url)
				}

				res, errSave := item.SaveToDB(db)

				if config.Debug {
					log.Printf("End save to DB %s \n", website.Url)
				}

				if errSave != nil {
					log.Print(errSave.Error())
					continue
				}
				if !res {
					continue
				}

				if config.Debug {
					log.Printf("News %s %s\n", item.Title, item.Url)
				}
				errSend := item.SaveToTg(config)
				if errSend != nil {
					log.Print(errSend.Error())
				}
				time.Sleep(5 * time.Second)
			}

			log.Printf("%s parsed \n", website.Url)
			errUpd := website.Update(db)
			if errUpd != nil {
				panic(errUpd)
			}

		} else {
			log.Printf("nothing to parse \n")
		}
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}
