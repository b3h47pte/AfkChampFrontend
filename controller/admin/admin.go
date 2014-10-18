/*
 * 'admin' handles the main admin page and general admin control related activities here.
 */
package admin

import(
  "net/http"
)

// Require the user to re-login. 
func RequireAdminRelogin(w http.ResponseWriter, r *http.Request) error {
  http.Redirect(w, r, "login?admin=true", 302)
  return nil
}

// Handles the home page for the admin section. 
func HandleAdminRoute(w http.ResponseWriter, r *http.Request) {
  if err := RequireAdminRelogin(w, r); err != nil {
    return
  }
  
  // If we get here then we know that the user is an admin.
}