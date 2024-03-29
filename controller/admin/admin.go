/*
 * 'admin' handles the main admin page and general admin control related activities here.
 */
package admin

import (
	"AfkChampFrontend/controller"
	"net/http"
)

type AdminTemplateData struct {
	Data controller.BaseTemplateData
}

const DefaultPageSize = 15

// 'RequireAdminRelogin' requires the user to re-login when visiting the admin page. Returns nil if user is already logged in as an admin.
func RequireAdminRelogin(w http.ResponseWriter, r *http.Request) error {
	if _, err := controller.GetCurrentUser(w, r, true); err != nil {
		http.Redirect(w, r, "/login?admin=true", 302)
		return err
	}
	return nil
}

// 'HandleAdminRoute' Handles the home page for the admin section.
func HandleAdminRoute(w http.ResponseWriter, r *http.Request) {
	if err := RequireAdminRelogin(w, r); err != nil {
		return
	}

	// If we get here then we know that the user is an admin.
	t := AdminTemplateData{Data: controller.CreateTemplateData(w, r)}
	controller.TemplateMapping["admin/admin.html"].ExecuteTemplate(w, "tbase", t)
}
