package main

import (
  "fmt"
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

func wait() {
  fmt.Println("\n\tEnter to quit")
  var q string
  fmt.Scanf("%s", &q)
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
  mode := <-web.Cmode
  if mode == "server" {
    go server.AsServer()
    wait()
  } 
  if mode == "text" {
    go func() {
      fmt.Println("\n\t TEXT")
      err := client.AsClient("text")
      if err != nil {
        fmt.Println(err)
      }
    }()
    wait()
  }
  if mode == "files" {
    go func() {
      fmt.Println("\n\t FILES")
      err := client.AsClient("files")
      if err != nil {
        fmt.Println(err)
      }
    }()
    wait()
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
}