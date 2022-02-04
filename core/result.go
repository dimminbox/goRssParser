package core

import (
	"database/sql"
	"errors"
	"fmt"
	"goRssParser/db"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Result struct {
	DateCreated    time.Time
	DatePublicated time.Time
	Url            string
	Hash           string
	Title          string
	Keyword        string
}

func (r Result) SaveToDB(db *sql.DB) (bool, error) {

	var sql string

	sql = fmt.Sprintf("select count(*) as cnt from results where url =\"%s\"", r.Url)
	rows, err := db.Query(sql)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	cnt := 0
	rows.Next()
	rows.Scan(&cnt)

	if cnt == 0 {
		sql = fmt.Sprintf("insert into results(date_created, url, title) values(\"%s\",\"%s\", \"%s\") ",
			r.DateCreated.Format("2006-01-02T15:04:05"), r.Url, strings.ReplaceAll(r.Title, "\"", ""))
		rows1, errI := db.Query(sql)
		defer rows1.Close()
		if errI != nil {
			return false, errI
		}
		return true, nil

	}
	return false, nil
}

func (r Result) SaveToTg(config db.Config) error {
	text := fmt.Sprintf("%s; %s; %s; %s", r.Title, r.Keyword, r.DatePublicated.Format("2006-01-02T15:04:05"), r.Url)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		config.Telegram.Token, config.Telegram.Chat, text)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return errors.New(fmt.Sprintf("error to send tg, status %d", err.Error()))
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("error to send tg, status %d, body %s", resp.StatusCode, string(body)))
	}
	return nil
}
