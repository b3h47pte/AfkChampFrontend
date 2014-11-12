/*
 * 'about' handles the display of the about page. This shooould be a static page but just in case.
 */
package about

import (
	"AfkChampFrontend/controller"
	"net/http"
)

type AboutTemplateData struct {
	Data controller.BaseTemplateData
}

// HandleAboutRoute displays the about page.
func HandleAboutRoute(w http.ResponseWriter, r *http.Request) {
	t := AboutTemplateData{Data: controller.CreateTemplateData(w, r)}
	controller.TemplateMapping["about/about.html"].ExecuteTemplate(w, "tbase", t)
}
