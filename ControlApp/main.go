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
	hardware, err := Infrastructure.Initialize()
	if err != nil {
		log.Fatalf("Error initializing hardware: %s", err)
	}

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

	//Handle lighting override endpoints
	http.HandleFunc("/api/lighting/auto", fixture.HandleSetLightingOverrideAutoApi)
	http.HandleFunc("/api/lighting/off", fixture.HandleSetLightingOverrideOffApi)
	http.HandleFunc("/api/lighting/static", fixture.HandleSetLightingOverrideSetColorApi)
	http.HandleFunc("/api/lighting/fade-to-static", fixture.HandleSetLightingOverrideFadeToColorApi)
	http.HandleFunc("/api/lighting/palette-fade", fixture.HandleSetLightingOverridePaletteFadeApi)
	http.HandleFunc("/api/lighting/palette-switch", fixture.HandleSetLightingOverridePaletteSwitchApi)
	http.HandleFunc("/api/lighting/brightness-flash", fixture.HandleSetLightingOverridePaletteBrightnessFlashApi)
	http.HandleFunc("/api/lighting/hue-flash", fixture.HandleSetLightingOverridePaletteHueFlashApi)
	http.HandleFunc("/api/lighting/strobe", fixture.HandleSetLightingOverrideStrobeApi)

	//Handle screen override endpoints
	http.HandleFunc("/api/screen/animation/override", fixture.HandleSetScreenOverrideAnimationSetApi)
	http.HandleFunc("/api/screen/text/override", fixture.HandleSetScreenOverrideTextSetApi)
	http.HandleFunc("/api/screen/brightness/level", fixture.HandleSetScreenOverrideBrightnessLevelApi)

	//Handle debug endpoints
	http.HandleFunc("/api/animations/import", fixture.HandleDisplayImportAnimationApi)
	http.HandleFunc("/api/display/connected", fixture.HandleDisplayConnectedApi)
	http.HandleFunc("/api/display/show", fixture.HandleDisplayPlayAnimationApi)
	http.HandleFunc("/api/display/text", fixture.HandleDisplayShowTextApi)
	http.HandleFunc("/api/display/brightness", fixture.HandleDisplaySetBrightnessApi)

	// Start server (listening on localhost prevents firewall popup on Windows)
	log.Println("Listening started")
	log.Fatalln(http.ListenAndServe("192.168.4.1:8080", nil))
}
