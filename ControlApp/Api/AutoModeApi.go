package Api

import (
	"ControlApp/BoxiBus"
	"ControlApp/Lightshow"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (fixture Fixture) HandleAutoModeSettingsApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fixture.handlePaletteGetApi(w, r)
	} else if r.Method == "POST" {
		fixture.handlePaletteCreateApi(w, r)
	} else if r.Method == "PUT" {
		fixture.handlePaletteUpdateApi(w, r)
	} else if r.Method == "DELETE" {
		fixture.handlePaletteDeleteApi(w, r)
	} else {

	}
}

func (fixture Fixture) handlePaletteGetApi(w http.ResponseWriter, r *http.Request) {
	var id uint32
	idStr := r.FormValue("id")
	if idStr != "" {
		tempId, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil || tempId < 0 {
			http.Error(w, "Error parsing ID.", http.StatusBadRequest)
			return
		}
		id = uint32(tempId)
	} else {
		http.Error(w, "ID not specified.", http.StatusBadRequest)
		return
	}

	exists, entity := fixture.Visuals.GetPalettes().GetById(id)

	if !exists {
		http.Error(w, "Palette does not exist.", http.StatusBadRequest)
		return
	}

	header := paletteHeader{entity.Id, entity.Name}
	var colors []color
	var moods []int

	for _, col := range entity.Colors {
		colors = append(colors, color{int(col.Red), int(col.Green), int(col.Blue), int(col.White), int(col.Amber), int(col.UltraViolet)})
	}

	for _, mood := range entity.Moods {
		moods = append(moods, int(mood))
	}

	palette := paletteType{header, moods, colors}

	//Encode data
	if err := json.NewEncoder(w).Encode(palette); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (fixture Fixture) handlePaletteCreateApi(w http.ResponseWriter, r *http.Request) {
	id := rand.Uint32()
	for exists, _ := fixture.Visuals.GetPalettes().GetById(id); exists; {
		id = rand.Uint32()
	}

	var data paletteCreate

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	var moods []Lightshow.LightingMood

	for _, mood := range data.Moods {
		if mood < 0 || mood > 3 {
			http.Error(w, fmt.Sprintf("Illegal mood value '%d'.", mood), http.StatusBadRequest)
			return
		}

		moods = append(moods, Lightshow.LightingMood(mood))
	}

	palette := Lightshow.Palette{Id: id, Name: data.Name, Moods: moods, Colors: []BoxiBus.Color{{}}}
	fixture.Visuals.GetPalettes().SetPalette(palette)

	returnData := paletteCreated{id}

	//Encode data
	if err := json.NewEncoder(w).Encode(returnData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
