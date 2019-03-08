package main

import (
	"fmt"
	"log"
	"net/http"
	h "twitch/server/handlers"

	"github.com/gorilla/mux"
)

func main() {

	var host = ":8888"
	r := mux.NewRouter()
	for ruta, manejador := range h.Manejadores {
		r.HandleFunc(ruta, manejador)
		fmt.Println("Manejador cargado en " + ruta)
	}
	fmt.Println("Servidor iniciándose en " + host)
	log.Fatal(http.ListenAndServe(host, r))
}
