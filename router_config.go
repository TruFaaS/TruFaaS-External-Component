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

// Initialize initializes the router configuration
func (routerConfig *RouterConfig) Initialize(platform FaaSPlatform) {
	routerConfig.Platform = platform
	routerConfig.Router = mux.NewRouter().StrictSlash(true)
	routerConfig.initializeOpenFaaSRoutes()

}

// Run Starts the router
func (routerConfig *RouterConfig) Run() {
	log.Fatal(http.ListenAndServe(":8080", routerConfig.Router))
}

// Fission Routes
func (routerConfig *RouterConfig) initializeFissionRoutes() {
	routerConfig.Router.HandleFunc("/test", TestFission)

}

// OpenFaaS Routes
func (routerConfig *RouterConfig) initializeOpenFaaSRoutes() {

}
