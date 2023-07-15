package client

import (
   "net"
   "crypto/rand"
   "crypto/sha256"
   "math/big"
   "encoding/json"
   "fmt"
   "encoding/base64"
   "errors"
   
   "github.com/yurajp/bridge/ascod"
   "github.com/yurajp/bridge/symcod"
)

type PassMode struct {
  Password string
  Mode string
}

var Mode string

func AsClient(mode string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", Conf.Addr, Conf.Port))
	if err != nil {
		fmt.Println("\t NO CONNECTION!")
		return err
	}
  defer conn.Close()
  // generating and sending random number
	rnd, err := ascod.GenRandom()
	if err != nil {
		return fmt.Errorf("Cannot generate random", err)
	}
	_, err = conn.Write([]byte(rnd))
	if err != nil {
		return fmt.Errorf("Can't send random", err)
	}
	// getting response with publick key and encrypted number 
	bfResp := make([]byte, 512)
	n, err := conn.Read(bfResp[:])
	if err != nil {
		return fmt.Errorf("Can't get response", err)
	}
	var resp ascod.KeyResp
	er := json.Unmarshal(bfResp[:n], &resp)
	if er != nil {
		return fmt.Errorf("Can't unjson response", err)
	}
	var pub ascod.PubKey
	// decrypted number must match sended number
	if pb, ok := ascod.GetPubAndCheck(resp, rnd); ok {
	  pub = pb
	} else {
	  return errors.New("Invalid KeyResponse")
	}
	// generating and assymetric encryptyng the password for futher using in connection
	pass := ascod.GeneratePasswrd(9)
	encPw := ascod.ClEncodeString(pass, pub)
	// Symmetric encoding Mode. It will be also used for password checking by server.
	encMode:= symcod.SymEncode(Mode, pass)
	// inserting both into struct and sending to server
	passMd := PassMode{encPw, encMode}
	jsPm, err := json.Marshal(passMd)
   if err != nil {
		return fmt.Errorf("Cannot marshal passMode: %w", err)
	}
	_, err = conn.Write([]byte(jsPm))
	if err != nil {
	  return fmt.Errorf("Cannot write passMode: %w", err)
	}
	srvConfirm := make([]byte, 64)
	p, err := conn.Read(srvConfirm[:])
	if err != nil {
	  return fmt.Errorf("Cannot read confirmation: %w",err)
	}
	sCf := string(srvConfirm[:p])
  if !ascod.IsClConfirmed(sCf, pub) {
		 return errors.New("Password is not confirmed")
  }
  if mode == "text" {
	  err = SendText(conn, pass)
	  if err != nil {
	    return err
	  }
  } 
  if mode == "file" {
	  err = SendFiles(conn)
	  if err != nil {
	    return err
	  }
  }
  return nil
}
