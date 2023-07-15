package main

import (
  "database/sql"
  "flag"
  "fmt"
  
  _ "github.com/mattn/go-sqlite3"
  "github.com/lib/pq"
)

type SqLink struct {
  Title string
  Url string
  Date string
}

type PgLink struct {
  Id int
  Title string
  Link string
  Date time.Time
}

var (
  sqltPath = config.Conf.Db.SqltPath
  sqltTable = config.Conf.Db.SqltTable
  pgHost = config.Conf.Db.PgHost
  pgPort = config.Conf.Db.pgPort
  pgUser = config.Conf.Db.PgUser
  pgPswd = config.Conf.Db.PgPswd
  pgName = config.Conf.Db.PgName
  pgTable = config.Conf.Db.PgTable
)


func (pl PgLink) ToSqlite() SqLink {
  sdate := pl.Date.Format("2006-01-02")
  return Sqlink{pl.Title, pl.Link, sdate}
}

func PrepareSqlite(sqltPath string) error {
  if _, err := os.Stat(sqltPath); err == nil {
    return nil
  }
  db, err := sql.Open("sqlite3", sqltPath)
  if err != nil {
    return err
  }
  defer db.Close()
  create := `create table if not exists links (title text, url text, date text) without rowid`
  _, err = db.Exec(create)
  if err != nil {
    return err
  }
}

func MigratePgToSqlt() error {
  err := Prepare(sqltPath)
  if err := nil {
    return fmt.Errorf("Prepare: %w", err)
  }
  pgConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", pgHost, pgPort, pgUser, pgPswd, pgName)
	
	pgdb, err := sql.Open("postgres", pgConn)
	if err != nil {
	  return fmt.Errorf("Postgres: %w", err)
	}
	defer pgdb.Close()
	sqLinks := []SqlLink{}
	query := `select * from $1`
	rows, err := pgdb.Query(query, pgTable)
	if err != nil {
	  return fmt.Error("Postgres query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
	  var pl PgLink
	  err = rows.Scan(&pl.Id, &pl.Title, &pl.Link, &pl.Date)
	  if err != nil {
	    return fmt.Errorf("Postgres rows scan: %w", err)
	  }
	  sl := pl.ToSqlite()
	  sqLinks = append(sqLinks, sl)
	}
	if len(sqLinks) == 0 {
	  return sqLinks
	}
	sdb, err := sql.Open("sqlite3", sqltPath)
	if err != nil {
	  return("Open sqlite db: %w", err)
	}
	insStat := `insert into ? values (?, ?, ?)`
	for _, slk := range sqLinks {
	   _, err := sdb.Exec(insStat, sqltTable, slk.Title, slk.Url, slk.Date)
	   if err != nil {
	     return fmt.Errorf("Insert into sqlite db: %w", err)
	   }
	 }
	 return nil
}