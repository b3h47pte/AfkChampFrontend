/*
 * 'error' Provides functionality to display error pages.
 */
package controller

import (
	"net/http"
)

func Handle404Page(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("404 Page"))
}

func HandleUnknownError(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Something bad happened, sorry! Try again later!"))
}
