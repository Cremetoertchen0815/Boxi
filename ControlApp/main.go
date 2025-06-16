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
	log.Println("Starting application...")

	// Initialize hardware
	hardware := &Infrastructure.DebugStub{}
	//hardware, err := Infrastructure.Initialize()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Initialize lighting manager
	visuals := Lightshow.CreateVisualManager(hardware)

	// Setup static file server
	fileServer := http.FileServer(http.Dir("Frontend/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Setup views
	pages := Frontend.CreatePageProvider()
	http.HandleFunc("/", pages.HandleStartPage)

	// Setup api
	fixture := Api.Fixture{Hardware: hardware, Visuals: visuals}
	http.HandleFunc("/api/animations/import", fixture.HandleDisplayImportAnimationApi)
	http.HandleFunc("/api/display/connected", fixture.HandleDisplayConnectedApi)
	http.HandleFunc("/api/display/show", fixture.HandleDisplayPlayAnimationApi)
	http.HandleFunc("/api/display/text", fixture.HandleDisplayShowTextApi)
	http.HandleFunc("/api/display/brightness", fixture.HandleDisplaySetBrightnessApi)
	http.HandleFunc("/api/debug/beat", func(w http.ResponseWriter, r *http.Request) {
		hardware.BeatTriggered = true
	})

	// Start server (listening on localhost prevents firewall popup on Windows)
	log.Println("Listening started")
	log.Fatalln(http.ListenAndServe("localhost:8080", nil))
}
