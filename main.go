package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strconv"

	webview "github.com/webview/webview_go"
)

//go:embed dist
var distFS embed.FS

func someApi(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Hello World from golang!")
}

func main() {
	fSys, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(fSys)))
	http.HandleFunc("/api/example", someApi)

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	go http.Serve(listener, nil)

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Basic Example")
	w.SetSize(1000, 618, webview.HintNone)
	w.Navigate("http://localhost:" + strconv.Itoa(port))
	w.Run()
}
