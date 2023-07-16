package web

import (
  "net/http"
  "html/template"
  "os/exec"
  "fmt"
  "github.com/yurajp/bridge/config"
  "github.com/yurajp/bridge/client"
)

var (
  hmTmpl *template.Template
  srTempl *template.Template
  clTempl *template.Template
  Cmode = make(chan string, 1)
)


func init() {
  fmt.Println("web inited")
  
  hm, err := template.ParseFiles("./files/hmTmpl.html")
  if err != nil {
    fmt.Println(err)
  } else {
    hmTmpl = hm
  }
  sr, err := template.ParseFiles("./files/srTmpl.html")
  if err != nil {
    fmt.Println(err)
  } else {
    srTmpl = sr
  }
  cl, err := template.ParseFiles("./files/clTmpl.html")
  if err != nil {
    fmt.Println(err)
  } else {
    clTmpl = cl
  }
}

func home(w http.ResponseWriter, r *http.Request) {
  hmTmpl.Execute(w, nil)
}

func serverLauncher(w http.ResponseWriter, r *http.Request) {
  Cmode <-"server"
  port := config.Conf.Server.Port
  serv := fmt.Sprintf("Bridge server is runing on %d", port)
  srTmpl.Execute(w, serv)
}

func textLauncher(w http.ResponseWriter, r *http.Request) {
  Cmode <-"text"
  for {
    select {
    case <-client.Res:
      clTmpl.Execute(w, client.Result)
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
      default:
    }
  }
}

func Launcher() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", home)
  mux.HandleFunc("/server", serverLauncher)
  mux.HandleFunc("/text", textLauncher)
  mux.HandleFunc("/files", filesLauncher)
  fs := http.FileServer(http.Dir("./files"))
  mux.Handle("/files/", http.StripPrefix("/files/", fs))
  hsrv := &http.Server{Addr: ":8642", Handler: mux}
  
  go hsrv.ListenAndServe()
  exec.Command("xdg-open", "http://localhost:8642").Run()
}