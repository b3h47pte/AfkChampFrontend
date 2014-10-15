package main

import(
  "github.com/gorilla/mux"
  "net/http"
  "AfkChampFrontend/controller/admin"
  "AfkChampFrontend/controller"
)

type AfkChampHandler func(w http.ResponseWriter, req *http.Request)

func main() {
  // Setup the routing for the frontend.
  r := mux.NewRouter()
  
  // Dynamic Content
  r.HandleFunc("/",controller.HandleHomeRoute)
  r.HandleFunc("/admin",admin.HandleAdminRoute)
  r.HandleFunc("/login",controller.HandleLoginPageRoute).Methods("GET")
  r.HandleFunc("/login",controller.HandleLoginAction).Methods("POST")
  r.HandleFunc("/register",controller.HandleRegisterPageRoute).Methods("GET")
  r.HandleFunc("/register",controller.HandleRegisterAction).Methods("POST")
  r.HandleFunc("/logout",controller.HandleLogoutPageRoute).Methods("GET")
  
  // Static Content
  r.PathPrefix("/javascript/").Handler(http.StripPrefix("/javascript/", http.FileServer(http.Dir("./javascript/"))))
  r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
  
  http.Handle("/",r)
  
  // TODO: Use HTTPS
  http.ListenAndServe("127.0.0.1:80",nil)
}