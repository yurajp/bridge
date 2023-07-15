package client

import (
   "fmt"
   "os/exec"
   "strconv"
 //  mand  "math/rand"
   "bufio"
   "net"
   "os"
   "regexp"
   "encoding/json"
   "errors"
   "time"
)

type Letter struct {
	Text string
	SHA  string
}

type Upfile struct {
	Fname string 
	Fsize int
	Isdir bool
}

var (
  tfile = config.Conf.Client.TxtFile
  fdir = config.Conf.Client.FileDir
)

func validText(txt string) bool {
   re := regexp.MustCompile(`\w`)
   return re.Match([]byte(txt))
}

func sendText(conn net.Conn, pw string) error {
  
	file, err := os.Open(tfile)
	if err != nil {
		return fmt.Errorf("Cannot open textfile to read: %w", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	txtBuf := make([]byte, 1024*96)
	ln, err := reader.Read(txtBuf[:])
	if err != nil {
		return fmt.Errorf("Cannot read textfile: %w", err)
	}
	text := string(msg[:ln])
	if !validText(text) {
		return errors.New("Empty letter")
	}
	encText := Symcode(text, pw)
	sha := HashStr(text)
	letter := Letter{encText, sha}
	jsl, err := json.Marshal(letter)
	if err != nil {
		return fmt.Errorf("Cannot marshal letter: %w", err)
	} 
	_, err = conn.Write(jsl)
	if err != nil {
		return fmt.Errorf("Cannot send letter: %w", err)
	}
	
	resBuf := make([]byte, 512)
	m, err := conn.Read(resBuf[:])
	if err != nil {
	  return fmt.Errorf("Cannot read result: %w", err)
	}
	res := string(resBuf[:m])
	if res != "OK" {
	  return errors.New("Server: " + res)
	}
	fmt.Println("\t The letter's received")
	
	err = LinkScanner(text)
	if err != nil {
	  return err
	}
	fmt.Println("\t Links are stored")
	
	return nil
}


func sendFiles(conn net.Conn) error {
  srvOk := func() (string, bool) {
    buf := make([]byte, 128) 
    n, err := conn.Read(buf[:])
    if err != nil {
      fmt.Printf("Read error: %s\n", err)
      return false
    }
    if m := string(buf[:n]); m != "ok" {
      fmt.Printf("Server error: %s\n", m)
      return m, false
    }
    return "", true
  }
  
  send := func(m string) {
    conn.Write([]byte(m))
  }
  
  upfiles := []Upfile{}
  files, err := os.ReadDir(fdir) 
  if err != nil {
    return fmt.Errorf("Cannot read files directory: %w", err)
  }
  if len(files) == 0 {
    return errors.New("Nothing to upload")
  }
	var n, sz int
	for _, f := range files {
	   name := f.Name()
	   inf, err := os.Stat(fdir + "/" + name)
	   if err != nil {
	      return fmt.Errorf("Cannot get Info", err)
	   }
	   size := int(inf.Size())
	   dir := false
		 if f.IsDir() {
			 s, err := makeZip(name)
			 if err != nil {
			   return err
			 }
	     size = s
	  	 name += ".zip"
	  	 dir = true
		}
		
		u := Upfile{name, size, dir}
		upfiles = append(upfiles, u)
		n++
		sz += size
	}
  sx := "s"
  if n == 1 {
    sx = ""
  }
  fmt.Printf("\n  uploading %d file%s, total size %s\n", n, sx, anyBytes(sz))
    
  for i, uf := range upfiles {
    juf, err := json.Marshal(uf)
    if err != nil {
      return fmt.Errorf("Error when marshal upfile: %w", err)
    }
    _, err = conn.Write(juf)
    if err != nil {
      return fmt.Errorf("Error when send upfile: %w", err)
    }
    msg, ok := srvOk()
    if !ok {
      conn.Close()
      return errors.New("Server error: ", msg)
    }
    name := fdir + "/" + uf.Name
    data, err := os.ReadFile(name)
    if err != nil {
      return fmt.Errorf("File reading error: %w", err)
    }
    _, err := conn.Write(data)
    if err != nil {
      return fmt.Errorf("Cannot send data: %w", err)
    }
    msg, ok = srvOk()
    if !ok {
      conn.Close()
      return errors.New("Server error: ", msg)
    }
    err = os.Remove(name)
    if err != nil {
      fmt.Println(err)
    }
    
    if i == n - 1 {
      send("finish")
      break
    }
    send("more")
    if msg, ok = srvOk(); !ok {
      return errors.New("Server error: ", msg)
    }
  }
  if msg, ok = srvOk(); !ok {
      return errors.New("Server error: ", msg)
  }
  fmt.Println("\t ✔ SUCCESS")
  return nil
}

func makeZip(dir string) (int, error) {
   wdr, err := os.Getwd()
   if err != nil {
      return 0, fmt.Errorf("Cannot get working dir: %w", err)
   }
   err = os.Chdir(fdir)
   if err != nil {
      return 0, fmt.Errorf("Cannot cd to conf: %w", err)
   }
   name := dir + ".zip"
   cmd := exec.Command("zip", "-4", "-rm", name, dir)
   err = cmd.Run()
   if err != nil {
      return 0, fmt.Errorf("Cannot zip %s: %w", name, err)
   }
   inf, err := os.Stat(name)
   if err != nil {
      return 0, fmt.Errorf("Cannot get info for %s: %w", name, err)
   }
   err = os.Chdir(wdr)
   if err != nil {
      return 0, fmt.Errorf("Cannot cd back to working dir: %w", err)
   }
   return int(inf.Size()), nil
}
