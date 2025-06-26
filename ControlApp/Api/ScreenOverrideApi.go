package Api

import (
	"net/http"
)

func (fixture Fixture) HandleSetScreenOverrideAnimationEnableApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

}

func (fixture Fixture) HandleSetScreenOverrideAnimationSetApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

}

func (fixture Fixture) HandleSetScreenOverrideTextEnableApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

}

func (fixture Fixture) HandleSetScreenOverrideTextSetApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

}

func (fixture Fixture) HandleSetScreenOverrideBrightnessLevelApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

}

func (fixture Fixture) HandleSetScreenOverrideBrightnessPulseApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

}
