package main

import (
  "fmt"
  "os"
  "github.com/yurajp/bridge/config"
  "github.com/yurajp/bridge/database"
  "github.com/yurajp/bridge/server"
  "github.com/yurajp/bridge/client"
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
  args := os.Args
  if len(args) > 2 {
    fmt.Println("Wrong arguments")
    return
  }
  if len(args) == 1 {
    go server.AsServer()
    wait()
  } else {
    if args[1] == "text" {
      go func() {
        err := client.AsClient("text")
        if err != nil {
          fmt.Println(err)
        }
      }()
      wait()
    }
    if args[1] == "files" {
      go func() {
        err := client.AsClient("files")
        if err != nil {
          fmt.Println(err)
        }
      }()
      wait()
    }
    if args[1] == "config" {
      err := config.TerminalConfig()
      if err != nil {
        fmt.Println(err)
      }
    }
    if args[1] == "migrate" {
      err := database.MigratePgToSqlt()
      if err != nil {
        fmt.Println(err)
      }
    }
  }
}