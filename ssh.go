package main

import (
	"io/ioutil"
	"log"
	"runtime/debug"
	"time"

	"golang.org/x/crypto/ssh"
)

func NewSSHConfig(user, key string, timeout int, ignoreInsecure bool) (cfg *ssh.ClientConfig) {
	cfg = new(ssh.ClientConfig)
	cfg.User = user
	cfg.Auth = []ssh.AuthMethod{
		PrivateKey(key),
	}
	if ignoreInsecure {
		cfg.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}
	cfg.Timeout = time.Duration(time.Duration(timeout) * time.Second)
	return
}
func PrivateKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func (c *CMD) SetBuffers() {
	c.StdOut.Buffer = make(chan []byte, 1000000)
	c.StdErr.Buffer = make(chan []byte, 1000000)

	c.Session.Stdout = &c.StdOut
	c.Session.Stderr = &c.StdErr
	newSTDin, err := c.Session.StdinPipe()
	if err != nil {
		c.Session.Close()
		log.Println("STDOUT:", err)
		return
	}
	c.StdIn = newSTDin
}
func (c *CMD) SetBuffersAndOpenShell() {
	c.SetBuffers()
	// THE SHELL NEEDS TO BE LAST!
	err := c.Session.Shell()
	if err != nil {
		log.Println(err, string(debug.Stack()))
	}

}
func (c *CMD) NewSessionForCommand(conn *ssh.Client) {
	session, err := conn.NewSession()
	if err != nil {
		log.Println("Session error:", err)
		return
	}
	c.Session = session
	c.Conn = conn
}
