package main

import (
	"ControlApp/Api"
	"ControlApp/Frontend"
	"ControlApp/Infrastructure"
	"ControlApp/Lightshow"
	"log"
	"net/http"
)

func main() {

	//Initialize hardware
	hardware := Infrastructure.DebugStub(0)
	//hardware, err := Infrastructure.Initialize()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//Initialize lighting manager
	visuals := Lightshow.CreateVisualManager(hardware)

	//Setup static file server
	fileServer := http.FileServer(http.Dir("Frontend/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	//Setup views
	pages := Frontend.CreatePageProvider()
	http.HandleFunc("/", pages.HandleStartPage)

	//Setup api
	fixture := Api.Fixture{Hardware: hardware, Visuals: visuals}
	http.HandleFunc("/api/display/connected", fixture.HandleDisplayConnectedApi)
	http.HandleFunc("/api/display/import", fixture.HandleDisplayImportAnimationApi)
	http.HandleFunc("/api/display/show", fixture.HandleDisplayPlayAnimationApi)
	http.HandleFunc("/api/display/text", fixture.HandleDisplayShowTextApi)
	http.HandleFunc("/api/display/brightness", fixture.HandleDisplaySetBrightnessApi)

	//Start server (listening on localhost prevents firewall popup on Windows)
	log.Fatalln(http.ListenAndServe("192.168.4.1:8080", nil))
}
