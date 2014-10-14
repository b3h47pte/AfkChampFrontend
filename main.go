package main

import(
  "github.com/gorilla/mux"
  "net/http"
  "AfkChampFrontend/controller/home"
  "AfkChampFrontend/controller/admin"
  "AfkChampFrontend/controller"
)

type AfkChampHandler func(w http.ResponseWriter, req *http.Request)

func main() {
  // Setup the routing for the frontend.
  r := mux.NewRouter()
  r.HandleFunc("/",home.HandleHomeRoute)
  r.HandleFunc("/admin",admin.HandleAdminRoute)
  r.HandleFunc("/login",controller.HandleLoginPageRoute).Methods("GET")
  r.HandleFunc("/login",controller.HandleLoginAction).Methods("POST")
  r.HandleFunc("/register",controller.HandleRegisterPageRoute).Methods("GET")
  r.HandleFunc("/register",controller.HandleRegisterAction).Methods("POST")
  r.HandleFunc("/logout",controller.HandleLogoutPageRoute).Methods("GET")
  http.Handle("/",r)
  
  // TODO: Use HTTPS
  http.ListenAndServe("127.0.0.1:80",nil)
}