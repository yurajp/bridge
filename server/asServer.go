package server

import (
	"encoding/json"
	"fmt"
	"net"
	"github.com/yurajp/bridge/config"
  "github.com/yurajp/bridge/ascod"
  
)

type KeyResp struct {
	Rand string
	Pub  PubKey
}

type PassMode struct {
  Password string
  Mode string
}

var (
  port = config.Conf.Server.Port
)

func SecureHandle(conn net.Conn) {
  // getting client random
	rndBuf := make([]byte, 512)
	n, err := conn.Read(rndBuf[:])
	if err != nil {
		fmt.Printf("cannot read random: %s", err)
		return
	}
	rand := string(rndBuf[:n])
	// generating keys 
	pub, priv, err := ascod.GenerateKeys()
	if err != nil {
		fmt.Printf("cannot generate keys: %s", err)
		return
	}
	// create KeyResp for client
	kRs := ascod.NewKeyResp(rand, pub, priv)
	// sending KeyResp json
	js, err := json.Marshal(kRs)
	if err != nil {
		fmt.Printf("cannot convert keyResp: %s", err)
		return
	}
	_, er := conn.Write(js)
	if er != nil {
		fmt.Printf("cannot send keyResp: %s", er)
		return
	}
	// getting struct with encrypted password for symmetric encoding and /
	// mode (files|text) encrypted by this password
	passMdBuf := make([]byte, 1024)
	m, err := conn.Read(passMdBuf[:])
	if err != nil {
		fmt.Printf("cannot receive passMode: %s", err)
		return
	}
	var encPM PassMode
	err = json.Unmarshal(passBuf[:m], &encPM)
	// handling the struct and getting password and mode
	decPwd := ascod.SrvDecodeString(encPM.Password, priv)
	mode := symcod.SymDecode(encPM.Mode, decPwd)
  sOk := ascod.SrvEncodeString("OK", priv)
  // define further action
	if mode == "file" {
	  conn.Write([]byte(sOk))
    go getFiles(conn)
	} else if mode == "text" {
    conn.Write([]byte(sOk))
	  go getText(conn, decPwd)
	} else {
	  conn.Write([]byte("error"))
	  conn.Close()
	}
}	
	
func AsServer() { // (stop chan struct{}) {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf(fmt.Errorf("Failed to establish connection: %s\n", err))
	}
	fmt.Print("\n\tSERVER started on :4545...\n ")
	for {
	  // select {
	  // case <-stop:  
	  //   break
	  // default:
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("connection failed ...continue...\n")
			continue
		}
		go SecureHandle(conn)
	}
}
