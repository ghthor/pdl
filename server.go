package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	config := struct {
		Addr, RedirectPort, SslPort string
		SslCert, SslKey             string
	}{
		"127.0.0.1", "8080", "8081",
		"tls-gen/cert.pem", "tls-gen/key.pem",
	}

	listenAddrHTTP := fmt.Sprintf("%s:%s", config.Addr, config.RedirectPort)
	listenAddrHTTPS := fmt.Sprintf("%s:%s", config.Addr, config.SslPort)

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

	err := http.ListenAndServeTLS(listenAddrHTTPS, config.SslCert, config.SslKey, handler)
	if err != nil {
		log.Fatal(err)
	}
}
