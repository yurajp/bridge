package main

import (
  "fmt"
  "time"
  "github.com/yurajp/bridge/config"
  "github.com/yurajp/bridge/database"
  "github.com/yurajp/bridge/server"
  "github.com/yurajp/bridge/client"
  "github.com/yurajp/bridge/web"
)

func iserr(err error) bool {
  if err != nil {
    fmt.Println(err)
    return true
  }
  return false
}

func main() {
  err := config.LoadConf()
  if iserr(err) {
    return
  }
  err = database.PrepareDb()
  if iserr(err) {
    return
  }
  go web.Launcher()
  Main:
  for {
    select {
    case mode := <-web.Cmode:
      if mode == "server" {
        if web.SrvUp {
          fmt.Println("  Server already running")
        } else {
          go server.AsServer()
        }
      } 
      if mode == "text" {
        go func() {
          fmt.Println("\n\t TEXT")
          err := client.AsClient("text")
          if err != nil {
            fmt.Println(err)
          }
        }()
      }
      if mode == "files" {
        go func() {
          fmt.Println("\n\t FILES")
          err := client.AsClient("files")
          if err != nil {
            fmt.Println(err)
          }
        }()
      }
      if mode == "config" {
        err := config.TerminalConfig()
        if err != nil {
          fmt.Println(err)
        }
      }
      if mode == "migrate" {
        err := database.MigratePgToSqlt()
        if err != nil {
          fmt.Println(err)
        }
      }
    case <-web.Q:
      break Main
  //
    }
  }
  fmt.Println("\t CLOSED")
  time.Sleep(3 * time.Second)
}