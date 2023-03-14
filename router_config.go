package main

import (
	"fmt"
	fission "github.com/TruFaaS/TruFaaS/fission"
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
	fmt.Println("Initializing Router")
	routerConfig.Platform = platform
	routerConfig.Router = mux.NewRouter().StrictSlash(true)
	routerConfig.initializeFissionRoutes()

}

// Run Starts the router
func (routerConfig *RouterConfig) Run() {
	log.Fatal(http.ListenAndServe(":8080", routerConfig.Router))
}

// Fission Routes
func (routerConfig *RouterConfig) initializeFissionRoutes() {
	fmt.Println("Initializing fission Routes")
	routerConfig.Router.HandleFunc("/fn/create", fission.Create).Methods(http.MethodPost)

}

// OpenFaaS Routes
func (routerConfig *RouterConfig) initializeOpenFaaSRoutes() {

}
