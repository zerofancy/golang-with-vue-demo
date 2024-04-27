package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	webview "github.com/webview/webview_go"
)

//go:embed dist
var distFS embed.FS

func checkPort(port int) bool {
	// 尝试连接本地的指定端口
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return true
	}
	defer conn.Close()
	return false
}

func someApi(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Hello World from golang!")
}

func main() {
	isDebug := !checkPort(3000)

	if isDebug {
		log.Default().Println("Debug Mode!")
		// debug下请求转发到前端页面
		targetURL, err := url.Parse("http://localhost:3000")
		if err != nil {
			log.Fatal(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
	} else {
		log.Default().Println("Release Mode!")
		// release下使用打包的网页内容
		fSys, err := fs.Sub(distFS, "dist")
		if err != nil {
			panic(err)
		}
		http.Handle("/", http.FileServer(http.FS(fSys)))
	}

	http.HandleFunc("/api/example", someApi)

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	go http.Serve(listener, nil)

	serveUrl := "http://localhost:" + strconv.Itoa(port)
	if isDebug {
		log.Default().Println("Serve Url is", serveUrl)
	}

	w := webview.New(false)
	defer w.Destroy()
	if isDebug {
		w.SetTitle("Basic Example - DEBUG")
	} else {
		w.SetTitle("Basic Example")
	}
	w.SetSize(1000, 618, webview.HintNone)
	w.Navigate(serveUrl)
	w.Run()
}
