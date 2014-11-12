/*
 * Displays the home page and whatever information that needs to be on it.
 */
package controller

import (
	"net/http"
)

type HomeTemplateData struct {
	Data BaseTemplateData
}

// HandleHomeRoute displays the main page.
func HandleHomeRoute(w http.ResponseWriter, r *http.Request) {
	t := HomeTemplateData{Data: CreateTemplateData(w, r)}
	TemplateMapping["main/home.html"].ExecuteTemplate(w, "tbase", t)
}
