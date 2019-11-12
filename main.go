package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/alecthomas/kingpin.v2"
)

type config struct {
	port string
	host string
}

var cfg = config{}

func setConfig() {
	kingpin.HelpFlag.Short('h')
	kingpin.Flag("port", "Target Port").Short('p').Default(":3000").StringVar(&cfg.port)
	kingpin.Flag("host", "Hostname").Default("localhost").StringVar(&cfg.host)
	kingpin.Parse()
}

func getLanding(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Bucket Landing</h1>")
}

func main() {
	setConfig()
	app := mux.NewRouter()

	app.HandleFunc("/", getLanding)
	log.Println("Server Running On https://" + cfg.host + cfg.port)
	http.ListenAndServe(cfg.port, app)
}
