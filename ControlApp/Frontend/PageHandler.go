package Frontend

import (
	"fmt"
	"net/http"
)

type startPageInformation struct {
	ScaffoldInformation
	stuff string
}

func (Me PageProvider) HandleStartPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	//Create data structure
	startData := startPageInformation{scaffoldData, "lol"}

	//Execute template
	err := Me.startPage.Execute(w, startData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}
