package main

import (
	"fmt"
	"github.com/TruFaaS/TruFaaS/constants"
	"github.com/TruFaaS/TruFaaS/fission"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type RouterConfig struct {
	Router   *mux.Router
	Platform constants.FaaSPlatform
}

// Initialize initializes the router configuration
func (routerConfig *RouterConfig) Initialize(platform constants.FaaSPlatform) {
	fmt.Println("Initializing Router")
	routerConfig.Platform = platform
	routerConfig.Router = mux.NewRouter().StrictSlash(true)
	routerConfig.initializeSpecifiedPlatformRoutes()

}

// Run Starts the router
func (routerConfig *RouterConfig) Run() {
	log.Fatal(http.ListenAndServe(":8080", routerConfig.Router))
}

// To initialize only specified FaaS platform routes
func (routerConfig *RouterConfig) initializeSpecifiedPlatformRoutes() {
	platform := routerConfig.Platform
	switch platform {
	case constants.Fission:
		routerConfig.initializeFissionRoutes()
	case constants.OpenFaaS:
		routerConfig.initializeOpenFaaSRoutes()
	default:
		routerConfig.initializeFissionRoutes()

	}

}

// Fission Routes
func (routerConfig *RouterConfig) initializeFissionRoutes() {
	fmt.Println("Initializing fission Routes")
	routerConfig.Router.HandleFunc("/fn/create", fission.CreateFnTrustValue).Methods(http.MethodPost)
	routerConfig.Router.HandleFunc("/fn/verify", fission.VerifyFnTrustValue).Methods(http.MethodPost)

}

// OpenFaaS Routes
func (routerConfig *RouterConfig) initializeOpenFaaSRoutes() {
	fmt.Println("Initializing OpenFaaS Routes")

}
