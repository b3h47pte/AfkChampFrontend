/*
 * Handles the admin pages.
 */
package admin

import(
  "net/http"
)

func HandleAdminRoute(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello world!!"))
}