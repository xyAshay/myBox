package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (info *serverInfo) getLanding(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

func (info *serverInfo) getPath(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	currPath := filepath.Join(info.root, path)
	log.Println("Current Path ", currPath)
	pathInfo, err := os.Stat(currPath)
	if err == nil && pathInfo.IsDir() {
		http.ServeFile(w, r, "./web/serve.html")
	} else {
		http.ServeFile(w, r, filepath.Join(info.root, path))

	}
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
			newFile.Size = "-"
		} else {
			newFile.Type = "file"
			newFile.Size = fmt.Sprintf("%.2f KB", float64(file.Size())/1024)
		}
		list = append(list, newFile)
	}
	data, _ := json.Marshal(list)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (info *serverInfo) downloadFile(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	w.Header().Set("Content-Disposition", "attachment")
	http.ServeFile(w, r, filepath.Join(info.root, path))
}

func (info *serverInfo) uploadFile(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	uploadPath := filepath.Join(info.root, path)
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("target")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v Bytes\n", handler.Size)

	tmpFile, err := os.Create(fmt.Sprintf("%s/%s", uploadPath, handler.Filename))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tmpFile.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	tmpFile.Write(data)
	fmt.Println("File Sucessfully Uploaded")
}

func (info *serverInfo) removeFile(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	filePath := filepath.Join(info.root, path)
	fmt.Println(filePath + "  Successfully Deleted ...")
	os.Remove(filePath)
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
