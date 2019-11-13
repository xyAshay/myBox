package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	port = kingpin.Flag("port", "Target Port").Short('p').Default(":3000").String()
	dev  = kingpin.Flag("dev", "Enable Developer Mode").Short('d').Bool()
)

type serverInfo struct {
	root   string
	router *mux.Router
}

func createServer(root string) *serverInfo {
	router := mux.NewRouter()
	if root == "" {
		root = "."
	}
	init := &serverInfo{
		root:   root,
		router: router,
	}
	log.Printf("Root : %s\n", init.root)
	router.HandleFunc("/", init.getLanding)
	router.HandleFunc("/serve/{path:.*}", init.getPath)
	return init
}

func (info *serverInfo) getLanding(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Bucket Landing</h1> <br> <a href='/serve/'>Get Started</a>")
}

func (info *serverInfo) getPath(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	log.Println("Current Path ", info.root, path)
	http.ServeFile(w, r, filepath.Join(info.root, path))
}

func (info *serverInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	info.router.ServeHTTP(w, r)
}

func main() {
	kingpin.Parse()
	app := createServer("./")

	fmt.Printf("Server Running on https://localhost%s \nDeveloper : %t\n", *port, *dev)
	log.Fatal(http.ListenAndServe(*port, app))
}
