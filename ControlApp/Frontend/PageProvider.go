package Frontend

import (
	"html/template"
	"net/http"
)

type PageProvider struct {
	startPage *template.Template
}

type ScaffoldInformation struct {
	PageName  string
	PageTitle string
}

// CreatePageProvider loads all templates and returns a PageProvider
func CreatePageProvider() PageProvider {
	start := template.Must(template.ParseFiles("Frontend/template/scaffold.gohtml", "Frontend/template/start.gohtml"))
	return PageProvider{
		startPage: start}
}

func GetScaffoldData(r *http.Request) ScaffoldInformation {
	page := r.URL.Path[1:]

	return ScaffoldInformation{
		PageName:  page,
		PageTitle: GetPageTitle(page)}
}

func GetPageTitle(pageName string) string {
	return "Boxi App"
}
