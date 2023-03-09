package main

type FaaSPlatform int

// Declare related constants for each direction starting with index 1
const (
	Fission   FaaSPlatform = iota + 1 // EnumIndex = 1
	OpenWhisk                         // EnumIndex = 2
	OpenFaaS                          // EnumIndex = 3
)
