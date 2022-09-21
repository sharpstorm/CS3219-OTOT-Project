package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HTTPHandler httprouter.Handle

type HTTPServer interface {
	Start()
	AddAssetRoute(route string, assetPath string)
	AddStaticRoute(route string, assetPath string)
	Get(route string, handler HTTPHandler)
	Post(route string, handler HTTPHandler)
	Put(route string, handler HTTPHandler)
	Delete(route string, handler HTTPHandler)
}

type httpServer struct {
	port   uint16
	router *httprouter.Router
}

func CreateHTTPServer(port uint16) HTTPServer {
	router := httprouter.New()
	return &httpServer{
		port:   port,
		router: router,
	}
}

func (server *httpServer) Start() {
	fmt.Print("Starting Server")
	if err := http.ListenAndServe(":8081", server.router); err != nil {
		fmt.Print("Failed to start the HTTP server")
	}
}

func (server *httpServer) AddAssetRoute(route string, assetPath string) {
	server.router.ServeFiles(route, http.Dir(assetPath))
}

func (server *httpServer) AddStaticRoute(route string, assetPath string) {
	server.router.GET(route, func(resp http.ResponseWriter, req *http.Request, p httprouter.Params) {
		http.ServeFile(resp, req, assetPath)
	})
}

func (server *httpServer) Get(route string, handler HTTPHandler) {
	server.router.GET(route, httprouter.Handle(handler))
}

func (server *httpServer) Post(route string, handler HTTPHandler) {
	server.router.POST(route, httprouter.Handle(handler))
}

func (server *httpServer) Put(route string, handler HTTPHandler) {
	server.router.PUT(route, httprouter.Handle(handler))
}

func (server *httpServer) Delete(route string, handler HTTPHandler) {
	server.router.DELETE(route, httprouter.Handle(handler))
}
