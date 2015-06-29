package match

import (
	"AfkChampFrontend/controller"
	"net/http"
)

type MatchTemplateData struct {
	Data           controller.BaseTemplateData
}

/*
 * HandleMatchPageRoute displays the page for this particular match. 
 * Most of the interactive ability of the page is done via Javascript.
 */
func HandleMatchPageRoute(w http.ResponseWriter, r *http.Request) {
	t := createMatchMainTemplateData(w, r)
	controller.TemplateMapping["match/match.html"].ExecuteTemplate(w, "tbase", t)
}

func createMatchMainTemplateData(w http.ResponseWriter, r *http.Request) *MatchTemplateData {
	t := MatchTemplateData{Data: controller.CreateTemplateData(w, r)}
	return &t
}
