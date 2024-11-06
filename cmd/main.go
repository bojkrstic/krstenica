package main

import (
	"flag"
	"krstenica/krstenica-v1/api"
	"log"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config-file-path",
		"",
		"Path to file that contains configuration for users API")
}
func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	api.RunHTTPServer(configFilePath)
}
