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
	router.HandleFunc("/api/json/{path:.*}", init.getJSONlisting).Methods("GET")
	router.HandleFunc("/api/delete/{path:.*}", init.removeFile).Methods("DELETE")
	router.HandleFunc("/api/download/{path:.*}", init.downloadFile).Methods("GET")
	router.HandleFunc("/api/upload/{path:.*}", init.uploadFile).Methods("POST")
	router.HandleFunc("/api/assets/{path:.*}", init.getAssets).Methods("GET")
	router.HandleFunc("/serve/{path:.*}", init.getPath).Methods("GET")
	router.HandleFunc("/", init.getLanding).Methods("GET")
	return init
}

func main() {
	kingpin.Parse()
	app := createServer("./")

	fmt.Printf("Server Running on https://localhost%s \nDeveloper : %t\n", *port, *dev)
	log.Fatal(http.ListenAndServe(*port, app))
}
