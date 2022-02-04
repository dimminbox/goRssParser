package core

import (
	"database/sql"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	_ "github.com/ahmetb/go-linq/v3"
	"strings"
	"time"
)

type Keyword struct {
	Name  string
	Forms sql.NullString
}

type Website struct {
	Id       int
	Name     string
	Url      string
	Source   string
	Active   bool
	Priority int
	Timeout  int
}

func (w Website) Update(db *sql.DB) (errUpdate error) {
	var sql string
	sql = fmt.Sprintf("update websites set date_parsed = \"%s\" where id = %d", time.Now().Format("2006-01-02T15:04:05"), w.Id)
	rows, err := db.Query(sql)
	defer rows.Close()
	return err
}

func GetUrl(db *sql.DB, period int) (Website, error) {
	var sql string
	sql = fmt.Sprintf("select url, timeout, source, id from websites "+
		" where (UNIX_TIMESTAMP(date_parsed) + %d < UNIX_TIMESTAMP()) "+
		" OR (date_parsed IS NULL) order by priority limit 1", period)
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return Website{}, err
	}

	for rows.Next() {
		website := Website{}
		err2 := rows.Scan(&website.Url, &website.Timeout, &website.Source, &website.Id)
		if err2 != nil {
			fmt.Println(err2)
			return Website{}, err2
		}
		return website, nil
	}

	return Website{}, nil
}

func GetKeyword(db *sql.DB) ([]string, error) {

	keywordList := []string{}
	var sql string

	sql = "select name, forms from keywords;"
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {

		fmt.Println(err)
		return keywordList, err
	}

	for rows.Next() {

		keyword := Keyword{}
		err2 := rows.Scan(&keyword.Name, &keyword.Forms)
		if err2 != nil {
			fmt.Println(err2)
			return []string{}, err2
		}
		keywordList = append(keywordList, keyword.Name)
		if keyword.Forms.String != "" {
			form := strings.Split(keyword.Forms.String, ",")
			keywordList = append(keywordList, form...)
		}
	}

	result := []string{}
	linq.From(keywordList).Distinct().
		Select(func(c interface{}) interface{} {
			val := strings.Trim(c.(string), " ")
			return val
		}).
		ToSlice(&result)
	return result, nil
}

func GetKeywordExc(db *sql.DB) ([]string, error) {

	keywordList := []string{}
	var sql string

	sql = "select name, forms from keywordexceptions;"
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {

		fmt.Println(err)
		return keywordList, err
	}

	for rows.Next() {

		keyword := Keyword{}
		err2 := rows.Scan(&keyword.Name, &keyword.Forms)
		if err2 != nil {
			fmt.Println(err2)
			return []string{}, err2
		}
		keywordList = append(keywordList, keyword.Name)
		if keyword.Forms.String != "" {
			form := strings.Split(keyword.Forms.String, ",")
			keywordList = append(keywordList, form...)
		}
	}

	result := []string{}
	linq.From(keywordList).Distinct().
		Select(func(c interface{}) interface{} {
			val := strings.Trim(c.(string), " ")
			return val
		}).
		ToSlice(&result)
	return result, nil
}
