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

func init() {
  TemplateMapping = make(map[string]*template.Template)
  
  // Store a mapping of all templates here
  TemplateMapping["main/home.html"] = template.Must(template.New("main/home.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/main/home.html"))
  TemplateMapping["about/about.html"] = template.Must(template.New("about/about.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/about/about.html"))
  TemplateMapping["login/login.html"] = template.Must(template.New("login/login.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/login/login.html"))
  TemplateMapping["login/register.html"] = template.Must(template.New("login/register.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/login/register.html"))
}

func CreateTemplateData() BaseTemplateData {
  newTempData := BaseTemplateData{WebsiteName: "Rocket ELO"}
  return newTempData
}