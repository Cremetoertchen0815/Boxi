package main

import (
	"ControlApp/Api"
	"ControlApp/Frontend"
	"ControlApp/Logic"
	"log"
	"net/http"
)

func main() {

	//Initialize hardware
	hardware, err := Logic.InitializeHardware()
	if err != nil {
		log.Fatal(err)
	}

	//Setup static file server
	fileServer := http.FileServer(http.Dir("Frontend/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	//Setup views
	pages := Frontend.CreatePageProvider()
	http.HandleFunc("/", pages.HandleStartPage)

	//Setup api
	fixture := Api.Fixture{Hardware: hardware}
	http.HandleFunc("/api/display/connected", fixture.HandleDisplayConnectedApi)
	http.HandleFunc("/api/display/import", fixture.HandleDisplayImportAnimationApi)
	http.HandleFunc("/api/display/upload", fixture.HandleDisplayImportAnimationApi)
	http.HandleFunc("/api/display/show", fixture.HandleDisplayUploadAnimationApi)
	http.HandleFunc("/api/display/text", fixture.HandleDisplayShowTextApi)

	//Start server(listening on localhost prevents firewall popup on Windows)
	log.Fatalln(http.ListenAndServe("0.0.0.0:8080", nil))
}
