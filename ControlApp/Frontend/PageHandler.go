package Frontend

import (
	"fmt"
	"net/http"
)

type startPageInformation struct {
	ScaffoldInformation
	Mood              int
	Nsfw              bool
	Brightness        int
	ConnectedDisplays string
}

type overridePageInformation struct {
	ScaffoldInformation
	LightingMode int
}

func (Me PageProvider) HandleStartPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	//Create data structure
	mood := int(Me.Visuals.GetConfiguration().Mood)
	isNsfw := Me.Visuals.GetConfiguration().AllowNsfw
	brightness := int(Me.Visuals.GetBrightness() * 100)
	displays := fmt.Sprintf("%+v", Me.Hardware.GetConnectedDisplays())
	startData := startPageInformation{scaffoldData, mood, isNsfw, brightness, displays}

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.startPage.Execute(w, startData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (Me PageProvider) HandleOverridesPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)
	data := overridePageInformation{scaffoldData, 0}

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.overridesPage.Execute(w, data)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (Me PageProvider) HandleAnimationPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.animationsPage.Execute(w, scaffoldData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (Me PageProvider) HandlePalettesPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.palettesPage.Execute(w, scaffoldData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (Me PageProvider) HandleAutoPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.autoPage.Execute(w, scaffoldData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}
