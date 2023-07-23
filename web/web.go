package web

import (
  "net/http"
  "html/template"
  "os/exec"
  "fmt"
  "embed"
  "github.com/yurajp/bridge/config"
  "github.com/yurajp/bridge/client"
  "github.com/yurajp/bridge/server"
)

var (
  //go:embed files
  webDir embed.FS
  fs http.Handler
  hmTmpl *template.Template
  srTmpl *template.Template
  clTmpl *template.Template
  blTmpl *template.Template
  Cmode = make(chan string, 1)
  SrvUp bool
  Q = make(chan struct{}, 1)
)


func init() {
  fs = http.FileServer(http.FS(webDir))
  hmTmpl, _ = template.ParseFS(webDir, "files/hmTmpl.html")
  srTmpl, _ = template.ParseFS(webDir, "files/srTmpl.html")
  clTmpl, _ = template.ParseFS(webDir, "files/clTmpl.html")
  blTmpl, _ = template.ParseFS(webDir, "files/blank.html")
}

func home(w http.ResponseWriter, r *http.Request) {
  if SrvUp {
    port := config.Conf.Server.Port
    serv := fmt.Sprintf("server is runing on %s", port)
    srTmpl.Execute(w, serv)
  } else {
    hmTmpl.Execute(w, nil)
  }
}

func serverLauncher(w http.ResponseWriter, r *http.Request) {
  Cmode <-"server"
 // port := config.Conf.Server.Port
  srTmpl.Execute(w, server.ToWeb)
  SrvUp = true
}

func textLauncher(w http.ResponseWriter, r *http.Request) {
  Cmode <-"text"
  for {
    select {
      case <-client.Res:
      clTmpl.Execute(w, client.Result)
      return
      default:
    }
  }
}

func filesLauncher(w http.ResponseWriter, r *http.Request) {
  Cmode <-"files"
  for {
    select {
    case <-client.Res:
      clTmpl.Execute(w, client.Result)
      return
      default:
    }
  }
}


func quit(w http.ResponseWriter, r *http.Request) {
  err := blTmpl.Execute(w, "Bridge closed")
  if err != nil {
    fmt.Println(err)
  }
  Q <-struct{}{}
}

func Launcher() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", home)
  mux.HandleFunc("/server", serverLauncher)
  mux.HandleFunc("/text", textLauncher)
  mux.HandleFunc("/files", filesLauncher)
  mux.HandleFunc("/quit", quit)
  mux.Handle("/files/", fs)
  hsrv := &http.Server{Addr: ":8642", Handler: mux}
  
  go hsrv.ListenAndServe()
  exec.Command("xdg-open", "http://localhost:8642").Run()
}