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
  
  // Admin Templates
  TemplateMapping["admin/admin.html"] = template.Must(template.New("login/register.html").Delims("<<", ">>").ParseFiles("html/adminBase.html", "html/admin/admin.html"))
  
  // Admin Games
  TemplateMapping["admin/game/index.html"] = template.Must(template.New("admin/game/index.html").Delims("<<", ">>").ParseFiles("html/adminBase.html","html/admin/game/index.html"))
  TemplateMapping["admin/game/newedit.html"] = template.Must(template.New("admin/game/newedit.html").Delims("<<", ">>").ParseFiles("html/adminBase.html","html/admin/game/newedit.html"))
  
  // Admin Events
  
  // Admin Users
  TemplateMapping["admin/user/index.html"] = template.Must(template.New("admin/user/index.html").Delims("<<", ">>").ParseFiles("html/adminBase.html", "html/admin/user/index.html"))
}

func CreateTemplateData() BaseTemplateData {
  newTempData := BaseTemplateData{WebsiteName: "Rocket ELO"}
  return newTempData
}