package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

//FileListing : Generate File List JSON
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
	router.HandleFunc("/", init.getLanding)
	router.HandleFunc("/serve/{path:.*}", init.getPath)
	router.HandleFunc("/api/json/{path:.*}", init.getJSONlisting)
	return init
}

func (info *serverInfo) getLanding(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

func (info *serverInfo) getPath(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	log.Println("Current Path ", info.root, path)
	http.ServeFile(w, r, "./web/serve.html")
	//http.ServeFile(w, r, filepath.Join(info.root, path))
}

func (info *serverInfo) getJSONlisting(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	list := make([]FileListing, 0)
	fd, err := os.Open(filepath.Join(info.root, path))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer fd.Close()

	files, err := fd.Readdir(-1)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for _, file := range files {
		newFile := FileListing{
			Name: file.Name(),
			Path: filepath.Join(path, file.Name()),
		}
		if file.IsDir() {
			newFile.Type = "dir"
			newFile.Size = "-1"
		} else {
			newFile.Type = "file"
			newFile.Size = "0"
		}
		list = append(list, newFile)
	}
	data, _ := json.Marshal(list)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
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
