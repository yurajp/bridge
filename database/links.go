package database

import (
	"fmt"
	"os"
	"net/http"
	"regexp"
	"strings"
	"time"
	"bufio"
	"database/sql"
  "github.com/yurajp/bridge/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/PuerkitoBio/goquery"
)

type Link struct {
  Title string
  Url string
  Date string
}


var (
  dbFile = config.Conf.Db.SqltPath
  lkTable = config.Conf.Db.SqltTable
)

func PrepareDb() error {
  if _, err := os.Stat(dbFile); err == nil {
    return nil
  }
  db, err := sql.Open("sqlite3", dbFile)
  if err != nil {
    return err
  }
  defer db.Close()
  create := `create table if not exists ? (title text, url text, date text) without rowid`
  _, err = db.Exec(create, lkTable)
  if err != nil {
    return err
  }
  return nil
}

func ScrapeTitle(url string) string {
  cl := http.Client{Timeout: 10 * time.Second}
  resp, err := cl.Get(url)
  if err != nil {
    return ""
  }
  defer resp.Body.Close()
  if resp.StatusCode != 200 {
    return ""
  }
  doc, err := goquery.NewDocumentFromReader(resp.Body)
  if err != nil {
    return ""
  }
  var title string
	doc.Find("head").Each(func(_ int, s *goquery.Selection) {
		title = s.Find("title").Text()
	})
	return title
}

func LinkScanner(text string) error {
  url := regexp.MustCompile(`http(s)?://*`)
  sc := bufio.NewScanner(strings.NewReader(text))
  linksDb := []Link{}
  for sc.Scan() {
    line := sc.Text()
    if url.MatchString(line) {
      if title := ScrapeTitle(line); title != "" {
        link := Link{title, line, time.Now().Format("2006-01-02")}
        linksDb = append(linksDb, link)
      }
    }
  }
  if len(linksDb) == 0 {
    return nil
  }
  return handleDb(linksDb)
}


func handleDb(links []Link) error {
  err := PrepareDb()
  if err != nil {
    return fmt.Errorf("Error when creating table in db: %w", err)
  }
  db, err := sql.Open("sqlite3", dbFile)
  if err != nil {
    return fmt.Errorf("Cannot open database: %w", err)
  }
  defer db.Close()
  dedup := `DELETE FROM links WHERE link = ?`
  insert := `INSERT INTO links VALUES(?, ?, ?)`
  for _, lk := range links {
    _, err = db.Exec(dedup, lk.Url)
    if err != nil {
      return fmt.Errorf("Error when delete duplicate from db: %w", err)
    }
    _, err = db.Exec(insert, lk.Title, lk.Url, lk.Date)
    if err != nil {
      return fmt.Errorf("Error when insert into db: %w", err)
    }
  }
  return nil
}