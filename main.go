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

type serverInfo struct {
	root   string
	router *mux.Router
}

//FileListing : Generate File List for /api/JSON
type FileListing struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	Size string `json:"size"`
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
	router.HandleFunc("/api/json/{path:.*}", init.getJSONlisting)
	router.HandleFunc("/api/delete/{path:.*}", init.removeFile)
	router.HandleFunc("/api/download/{path:.*}", init.downloadFile)
	router.HandleFunc("/api/upload/{path:.*}", init.uploadFile)
	router.HandleFunc("/serve/{path:.*}", init.getPath)
	router.HandleFunc("/", init.getLanding)
	return init
}

func main() {
	kingpin.Parse()
	app := createServer("./")

	fmt.Printf("Server Running on https://localhost%s \nDeveloper : %t\n", *port, *dev)
	log.Fatal(http.ListenAndServe(*port, app))
}
