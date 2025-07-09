package Frontend

import (
	"ControlApp/Api"
	"html/template"
	"net/http"
)

type PageProvider struct {
	Data *Api.DataContainer

	//Templates
	startPage      *template.Template
	overridesPage  *template.Template
	animationsPage *template.Template
	palettesPage   *template.Template
	autoPage       *template.Template
}

type ScaffoldInformation struct {
	PageName  string
	PageTitle string
}

// CreatePageProvider loads all templates and returns a PageProvider
func CreatePageProvider(data *Api.DataContainer) PageProvider {
	//Load templates
	start := template.Must(template.ParseFiles("Frontend/template/scaffold.gohtml", "Frontend/template/start.gohtml"))
	overrides := template.Must(template.ParseFiles("Frontend/template/scaffold.gohtml", "Frontend/template/overrides.gohtml"))
	animations := template.Must(template.ParseFiles("Frontend/template/scaffold.gohtml", "Frontend/template/animations.gohtml"))
	palettes := template.Must(template.ParseFiles("Frontend/template/scaffold.gohtml", "Frontend/template/palettes.gohtml"))
	auto := template.Must(template.ParseFiles("Frontend/template/scaffold.gohtml", "Frontend/template/auto.gohtml"))

	return PageProvider{
		Data:           data,
		startPage:      start,
		overridesPage:  overrides,
		animationsPage: animations,
		palettesPage:   palettes,
		autoPage:       auto,
	}
}

func GetScaffoldData(r *http.Request) ScaffoldInformation {
	page := r.URL.Path[1:]

	return ScaffoldInformation{
		PageName:  page,
		PageTitle: GetPageTitle(page)}
}

func GetPageTitle(pageName string) string {
	switch pageName {
	case "overrides":
		return "Manual Lighting"
	case "animations":
		return "Animations"
	case "palettes":
		return "Palettes"
	case "auto":
		return "Auto Mode Settings"
	}

	return "Boxi Control App"
}
