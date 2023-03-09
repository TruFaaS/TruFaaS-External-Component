package main

func main() {
	routerConfig := RouterConfig{}
	routerConfig.Initialize(OpenFaaS)
	routerConfig.Run()
}
