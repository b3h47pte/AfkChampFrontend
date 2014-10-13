/*
 * 'View' Provides basic functionality for constructing a view using Go's template engine.
 */
package controller

import(
  "html/template"
)
var TemplateMapping map[string]*template.Template

type BaseTemplateData struct {
  WebsiteName string
  Data interface{}
}

func InitializeTemplates() {
  TemplateMapping = make(map[string]*template.Template)
  
  // Store a mapping of all templates here
  TemplateMapping["login/login.html"] = template.Must(template.ParseFiles("html/base.html", "html/login/login.html"))
  TemplateMapping["login/register.html"] = template.Must(template.ParseFiles("html/base.html", "html/login/register.html"))
}

func CreateTemplateData() BaseTemplateData {
  newTempData := BaseTemplateData{WebsiteName: "Rocket ELO"}
  return newTempData
}