package main

import (
	"WebSiteStatusChecker/Service"
	"log"
	"net/http"
)

func main() {
	//Request Mapping
	http.HandleFunc("/checkStatus", Service.HelloHandler)
	//Starting Server, Using Default Mux
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
