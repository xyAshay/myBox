package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	port = kingpin.Flag("port", "Target Port").Short('p').Default(":3000").String()
	dev  = kingpin.Flag("dev", "Enable Developer Mode").Short('d').Bool()
)

func getLanding(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Bucket Landing</h1>")
}

func main() {
	kingpin.Parse()
	app := mux.NewRouter()

	app.HandleFunc("/", getLanding)
	fmt.Printf("Server Running on https://localhost:%s \nDeveloper : %t", *port, *dev)
	log.Fatal(http.ListenAndServe(*port, app))
}
