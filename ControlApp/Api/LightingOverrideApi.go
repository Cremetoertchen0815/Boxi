package Api

import (
	"encoding/json"
	"net/http"
)

func (fixture Fixture) HandleSetLightingOverrideAutoApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	fixture.Visuals.SetLightingOverwrite(nil)

	//Encode data
	if err := json.NewEncoder(w).Encode(fetchResult{indices}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
