package main

import (
	"fmt"
	"github.com/ghthor/pdl/config"
	"log"
	"net/http"
)

func main() {
	log.Println("Reading Config file: config.json")

	config, err := config.ReadFromFile("config.json")
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	listenAddrHTTP := fmt.Sprintf("%s:%d", config.LAddr, config.RedirectPort)
	listenAddrHTTPS := fmt.Sprintf("%s:%d", config.LAddr, config.SslPort)

	addrHTTPS := fmt.Sprintf("https://%s", listenAddrHTTPS)

	go func() {
		log.Println("Starting...HTTP Redirect Server")

		server := &http.Server{
			Addr:    listenAddrHTTP,
			Handler: http.RedirectHandler(addrHTTPS, http.StatusMovedPermanently),
		}

		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Starting...HTTPS Web Server")

	handler := http.DefaultServeMux

	err = http.ListenAndServeTLS(listenAddrHTTPS, config.SslCert, config.SslKey, handler)
	if err != nil {
		log.Fatal(err)
	}
}
