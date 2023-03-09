package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type RouterConfig struct {
	Router   *mux.Router
	Platform FaaSPlatform
}

func (routerConfig *RouterConfig) Initialize(platform FaaSPlatform) {
	routerConfig.Platform = platform
	routerConfig.Router = mux.NewRouter().StrictSlash(true)
	routerConfig.initializeOpenFaaSRoutes()

}

func (routerConfig *RouterConfig) initializeFissionRoutes() {
	routerConfig.Router.HandleFunc("/test", TestFission)

}

func (routerConfig *RouterConfig) initializeOpenFaaSRoutes() {
	routerConfig.Router.HandleFunc("/test", TestOF)

}

func (routerConfig *RouterConfig) Run() {
	log.Fatal(http.ListenAndServe(":8080", routerConfig.Router))
}
