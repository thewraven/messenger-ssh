package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	user = flag.String("user", os.Getenv("SSH_CLIENT"), "usuario ssh -> $SSH_CLIENT")
	ruta = flag.String("archivo", os.Getenv("SSH_PUBLIC_KEY"), "archivo con llave pública $SSH_PUBLIC_KEY")
	//pass = flag.String("pass", os.Getenv("SSH_PASS"), "pass ssh -> $SSH_PASS")
)

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		panic(err)
	}
	//fmt.Println(key.PublicKey().Type())
	return ssh.PublicKeys(key)
}

func main() {
	flag.Parse()
	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			publicKeyFile(*ruta),
		},
	}
	conn, err := ssh.Dial("tcp", "localhost:22", config)
	checkErr(err)
	defer conn.Close()
	session, err := conn.NewSession()
	checkErr(err)
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	pipe, err := session.StdinPipe()
	defer pipe.Close()
	tee := io.TeeReader(os.Stdin, pipe)
	//Pipe entre Stdin local y Stdin de la sesión ssh
	leerDatos := func(r io.Reader) {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", b)
	}
	go leerDatos(tee)
	checkErr(err)

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	err = session.RequestPty("xterm", 120, 180, modes)
	checkErr(err)
	if err := session.Shell(); err != nil {
		panic(err)
	}
	err = session.Wait()
	fmt.Println("finalizando sesión con error ", err)
}
