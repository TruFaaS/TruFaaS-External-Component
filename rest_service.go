package main

func main() {
	routerConfig := RouterConfig{}
	routerConfig.Initialize(Fission)
	routerConfig.Run()
}
