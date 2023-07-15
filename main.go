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
  fmt.Println("   Enter to quit")
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
  }
  if args[1] == "text" {
    go client.AsClient("text")
  }
  if args[1] == "files" {
    go client.AsClient("files")
  }
  wait()
}