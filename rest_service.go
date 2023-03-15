package main

import "github.com/TruFaaS/TruFaaS/constants"

func main() {
	routerConfig := RouterConfig{}
	routerConfig.Initialize(constants.Fission)
	routerConfig.Run()
}
