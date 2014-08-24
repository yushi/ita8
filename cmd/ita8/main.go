package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/yushi/ita8"
)

func pbpaste() ([]byte, error) {
	c := exec.Command("pbpaste")
	return c.Output()
}

func pbcopy(in []byte) error {
	c := exec.Command("pbcopy")
	c.Stdin = bytes.NewBuffer(in)
	b, err := c.CombinedOutput()
	if err != nil {
		log.Println(string(b))
	}
	return err
}

func openCommand(args []string) error {
	c := exec.Command("open", args...)
	b, err := c.CombinedOutput()
	if err != nil {
		log.Println(string(b))
	}
	return err
}

func checkRemoteAddr(r *http.Request, allowed string) bool {
	addrs := strings.Split(":", r.RemoteAddr)
	return addrs[0] == allowed
}

func getClipboardHandler(remoteAddr string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if checkRemoteAddr(r, remoteAddr) {
			w.WriteHeader(401)
			return
		}

		switch r.Method {
		case "GET":
			b, err := pbpaste()
			if err != nil {
				w.WriteHeader(500)
			} else {
				w.WriteHeader(200)
			}
			w.Write(b)
		case "PUT":
			defer r.Body.Close()
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
			}
			err = pbcopy(b)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
			} else {
				w.WriteHeader(200)
			}
		default:
			w.WriteHeader(405)
		}
	}
}

func getOpenHandler(remoteAddr string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if checkRemoteAddr(r, remoteAddr) {
			w.WriteHeader(401)
			return
		}

		switch r.Method {
		case "POST":
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(400)
				return
			}
			v := &[]string{}
			err = json.Unmarshal(b, v)
			if err != nil {
				log.Println(err)
				w.WriteHeader(400)
				return
			}
			err = openCommand(*v)
			if err != nil {
				w.WriteHeader(500)
			} else {
				w.WriteHeader(200)
			}
		default:
			w.WriteHeader(405)
		}
	}
}

func getSSHCmd(remoteAddr, localAddr string) *exec.Cmd {
	cmdStr := fmt.Sprintf("bash -c 'killall ita8br; %s %s'", "ita8br", localAddr)
	return exec.Command("ssh", remoteAddr, cmdStr)
}

func checkConnectivity(remote string) (string, error) {
	conn, err := net.Dial("tcp", remote+":22")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := strings.Split(conn.LocalAddr().String(), ":")[0]
	return localAddr, nil
}

func main() {
	flag.Parse()
	remoteAddr := flag.Arg(0)
	if remoteAddr == "" {
		log.Fatal("remote addr required")
	}

	localAddr, err := checkConnectivity(remoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(localAddr)
	sshCmd := getSSHCmd(remoteAddr, localAddr)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		sshCmd.Process.Signal(syscall.SIGTERM)
	}()

	go func() {
		defer sshCmd.Process.Kill()
		b, err := sshCmd.CombinedOutput()
		if err != nil {
			log.Fatal(fmt.Sprintf("ssh error: %s", err.Error()))
		}
		log.Println(b)
		log.Fatal("ssh for ita8br closed.")
	}()

	http.HandleFunc(
		fmt.Sprintf("/%s", ita8.ClipboardPath),
		getClipboardHandler(remoteAddr),
	)
	http.HandleFunc(
		fmt.Sprintf("/%s", ita8.OpenPath),
		getOpenHandler(remoteAddr),
	)
	log.Fatal(http.ListenAndServe(":4567", nil))
}
