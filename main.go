package main

import(
  "github.com/gorilla/mux"
  "net/http"
  "AfkChampFrontend/controller/home"
  "AfkChampFrontend/controller/admin"
  "AfkChampFrontend/controller"
)

func main() {

  // Initialize frontend for use
  controller.InitializeTemplates()
  controller.InitializeLogin()

  // Setup the routing for the frontend.
  r := mux.NewRouter()
  r.HandleFunc("/",home.HandleHomeRoute)
  r.HandleFunc("/admin",admin.HandleAdminRoute)
  r.HandleFunc("/login",controller.HandleLoginPageRoute).Methods("GET")
  http.Handle("/",r)
  http.ListenAndServe("127.0.0.1:80",nil)
}