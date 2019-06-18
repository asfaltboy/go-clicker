package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
)

type LogRecord struct {
	Url string `json:"url"`
}

type Context struct {
	Host string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	fmt.Printf("Args: %s\n", os.Args)
	ngrok := flag.Bool("ngrok", false, "enable ngrok")
	hostname := flag.String("hostname", "localhost", "a hostname of the websocket server to connect to")
	port := flag.String("port", "8088", "a port of the websocket server to connect to")
	flag.Parse()

	var host string
	if !*ngrok {
		host = *hostname + ":" + *port
	} else {
		cmd := exec.Command(
			"ngrok", "tcp", "--region=eu", "--log-format=json", "--log=stdout", "8088",
		)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		var record LogRecord
		for {
			if err := json.NewDecoder(stdout).Decode(&record); err != nil {
				log.Fatal(err)
			}
			fmt.Print(".")
			if record.Url != "" {
				break
			}
		}
		fmt.Println("\n> parsing ngrok url", record.Url)
		parts := strings.Split(record.Url, "/")
		host = parts[len(parts)-1]
		fmt.Println("> setting host to:", host)
		fmt.Println("> Ready to receive", host)
	}

	http.HandleFunc("/click", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Fatal((err))
			}

			if string(msg) == "secret-token-12345" {
				fmt.Printf("%s sent click\n", conn.RemoteAddr())
				robotgo.MouseClick("left", true)
				err = conn.WriteMessage(msgType, []byte("OK"))
				if err != nil {
					log.Fatal((err))
				}
			}
		}
	})

	c := Context{Host: host}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("web/websockets.html")
		if err != nil {
			log.Fatal(err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, c)
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	http.ListenAndServe(":8088", nil)
}
