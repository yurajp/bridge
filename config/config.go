package config

import (
  "os"
  "encoding/json"
  "fmt"
)

type Config struct {
  Appdir string `json:"appdir"`
  Server Srv `json:"server"`
  Client Clt `json:"client"`
  Db Dbs `json:"db"`
}

type  Srv struct {
    Port string `json:"port"`
    TxtDir string `json:"txtdir"`
    FileDir string `json:"filedir"`
}
  
type  Clt struct {
    Addr string `json:"addr"`
    TxtFile string `json:"txtfile"`
    FileDir string `json:"filedir"`
}

type Dbs struct {
    SqltPath string `json:"sqltpath"`
    SqltTable string `json:"sqlttable"`
    PgHost string `json:"pghost"`
    PgPort string `json:"pgport"`
    PgUser string `json:"pguser"`
    PgPswd string `json:"pgpswd"`
    PgName string `json: "pgname"`
    PgTable string `json:"pgtable"`
}

var (
  STORAGE = "/storage/emulated/0/"
  TERMUXDIR = "/data/data/com.termux/files/home/"
)

var DefaultConf = Config{Appdir: TERMUXDIR + "golangs/bridge/mobile",
  Server: Srv{Port: ":4646",
    TxtDir: STORAGE + "BridgeTexts",
    FileDir: STORAGE + "BridgeFiles"},
  Client: Clt{Addr: "192.168.1.38:4545",
    TxtFile: STORAGE + "pc.txt",
    FileDir: STORAGE + "Uploads"},
  Db: Dbs{SqltPath: TERMUXDIR + "golangs/bridge-mobile/database/bridge.db",
    SqltTable: "links",
    PgHost: "localhost",
    PgPort: "5432",
    PgUser: "yura",
    PgPswd: "26335",
    PgName: "notedb",
    PgTable: "messer"}}
  
var Conf Config

func LoadConf() error {
  path := TERMUXDIR + "golangs/bridge-mobile/config/config.json"
  if _, err := os.Stat(path); err != nil {
    if os.IsNotExist(err) {
      er := TerminalConfig()
      if er != nil {
        return er
      }
    }
  }
  jf, err := os.ReadFile(path)
  if err != nil {
    return err
  }
  err = json.Unmarshal(jf, &Conf)
  if err != nil {
    return err
  }
  return nil
}

func WriteConf(cf Config) error {
  jf, err := json.MarshalIndent(cf, " ", "   ")
  if err != nil {
    return fmt.Errorf("Marshal conf: %w", err)
  }
  return os.WriteFile("config.json", jf, 0640)
}

func TerminalConfig() error {
  cf := Config{}
  set := func(f, dft string) string {
    fmt.Printf(" - %s (empty for default)\n", f)
    var s string
    fmt.Scanf("%s", &s)
    if s == "" {
      s = dft
    }
    return s
  }
  
  fmt.Println("\t APPDIR")
  cf.Appdir = set("application directory", DefaultConf.Appdir)
  fmt.Println("\t SERVER:")
  cf.Server.Port = set("port for server", DefaultConf.Server.Port)
  cf.Server.TxtDir = set("directory for texts", DefaultConf.Server.TxtDir)
  cf.Server.FileDir = set("directory for files", DefaultConf.Server.FileDir)
  fmt.Println("\t CLIENT:")
  cf.Client.Addr = set("remote server address", DefaultConf.Client.Addr)
  cf.Client.TxtFile = set("path to text file", DefaultConf.Client.TxtFile)
  cf.Client.FileDir = set("directory to upload", DefaultConf.Client.FileDir)
  fmt.Println("\t DATABASE:")
  cf.Db.SqltPath = set("path to Sqlite db", DefaultConf.Db.SqltPath)
  cf.Db.SqltTable = set("table name", DefaultConf.Db.SqltTable)
  fmt.Println(" Will you use Postgresql? [y/n]")
  var p string
  fmt.Scanf("%s", &p)
  if p == "y" {
    cf.Db.PgHost = set("pg host", DefaultConf.Db.PgHost)
    cf.Db.PgPort = set("pg port", DefaultConf.Db.PgPort)
    cf.Db.PgUser = set("pg user", DefaultConf.Db.PgUser)
    cf.Db.PgPswd = set("pg password", DefaultConf.Db.PgPswd)
    cf.Db.PgName = set("pg db's name", DefaultConf.Db.PgName)
    cf.Db.PgTable = set("table name", DefaultConf.Db.PgTable)
  }
  return WriteConf(cf)
}

