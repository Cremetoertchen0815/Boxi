package main

import (
	"ControlApp/Api"
	"ControlApp/Frontend"
	"ControlApp/Infrastructure"
	"ControlApp/Lightshow"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Starting application...")

	// Initialize hardware
	hardware := Infrastructure.DebugStub{}

	// Initialize lighting manager
	visuals := Lightshow.CreateVisualManager(hardware)

	// Setup static file server
	fileServer := http.FileServer(http.Dir("Frontend/template/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Setup views
	pages := Frontend.CreatePageProvider(hardware, visuals)
	http.HandleFunc("/", pages.HandleStartPage)
	http.HandleFunc("/auto", pages.HandleAutoPage)
	http.HandleFunc("/overrides", pages.HandleOverridesPage)
	http.HandleFunc("/palettes", pages.HandlePalettesPage)
	http.HandleFunc("/animations", pages.HandleAnimationPage)

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
	http.HandleFunc("/api/screen/animation", fixture.HandleSetScreenOverrideAnimationSetApi)
	http.HandleFunc("/api/screen/text", fixture.HandleSetScreenOverrideTextSetApi)
	http.HandleFunc("/api/screen/brightness", fixture.HandleSetScreenOverrideBrightnessLevelApi)
	http.HandleFunc("/api/screen/connected", fixture.HandleScreensConnectedApi)

	//Handle palette endpoints
	http.HandleFunc("/api/palettes", fixture.HandlePaletteGetAllApi)
	http.HandleFunc("/api/palette", fixture.HandleSinglePaletteApi)

	//Handle animation endpoints
	http.HandleFunc("/api/animations", fixture.HandleAnimationsGetAllApi)
	http.HandleFunc("/api/animation", fixture.HandleSingleAnimationApi)

	//Handle auto mode config endpoints
	http.HandleFunc("/api/config/mood", fixture.HandleChangeAutoModeMoodApi)
	http.HandleFunc("/api/config/nsfw", fixture.HandleChangeAutoModeNsfwApi)
	http.HandleFunc("/api/config/advanced", fixture.HandleChangeAutoModeConfigApi)

	//Handle other endpoints
	http.HandleFunc("/api/ping", func(writer http.ResponseWriter, request *http.Request) {})

	//Mark lightshow dirty after time delay
	go func(manager *Lightshow.VisualManager) {
		time.Sleep(time.Second * 5)
		manager.MarkLightshowAsDirty()
	}(visuals)

	// Start server (listening on localhost prevents firewall popup on Windows)
	log.Println("Listening started")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
