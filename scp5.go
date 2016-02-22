package main

import (
       "fmt"
       "flag"
       "io/ioutil"
       "log"
       "strings"
       "golang.org/x/crypto/ssh"
       "github.com/tmc/scp"
)

func split_connection_str(conn string) (user, host, path string, err error) {
    if strings.Index(conn, "@") < 0 {
        return "", "", "", fmt.Errorf("no @")
    }

    if strings.Index(conn, ":") < 0 {
        return "", "", "", fmt.Errorf("no :")
    }

    return conn[:strings.Index(conn, "@")], 
    	   conn[strings.Index(conn, "@")+1:strings.Index(conn, ":")],
	   conn[strings.Index(conn, ":")+1:],
	   nil
}

func main() {
    flag.Parse()

//     fmt.Println("arg num : ", len(flag.Args()))
//     fmt.Println("arg #1 : ", flag.Args()[0])

    args := flag.Args()

    user, host, path, err := split_connection_str(args[2])

    pemBytes, err := ioutil.ReadFile(args[0])
    if err != nil {
        log.Fatal(err)
    }
    signer, err := ssh.ParsePrivateKey(pemBytes)
    if err != nil {
        log.Fatalf("parse key failed:%v", err)
    }
    config := &ssh.ClientConfig{
        User: user,
	Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
    }
    conn, err := ssh.Dial("tcp", host + ":22", config)
    if err != nil {
        log.Fatalf("dial failed:%v", err)
    }
    defer conn.Close()
    session, err := conn.NewSession()
    if err != nil {
        log.Fatalf("session failed:%v", err)
    }
    defer session.Close()

    err = scp.CopyPath(args[1], path, session)
}