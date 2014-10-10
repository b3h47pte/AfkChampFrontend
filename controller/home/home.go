/*
 * Displays the home page and whatever information that needs to be on it.
 */
package home

import(
  "net/http"
)

func HandleHomeRoute(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello world"))
}