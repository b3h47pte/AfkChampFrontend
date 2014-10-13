package main

import(
  "github.com/gorilla/mux"
  "net/http"
  "AfkChampFrontend/controller/home"
  "AfkChampFrontend/controller/admin"
  "AfkChampFrontend/controller"
  "AfkChampFrontend/model"
  "AfkChampFrontend/model/user"
  "log"
)

type AfkChampHandler func(w http.ResponseWriter, req *http.Request)

func main() {
  model.InitializeDatabase()
  controller.InitializeTemplates()
  controller.InitializeLogin()

  // Setup the routing for the frontend.
  r := mux.NewRouter()
  r.HandleFunc("/",home.HandleHomeRoute)
  r.HandleFunc("/admin",admin.HandleAdminRoute)
  r.HandleFunc("/login",controller.HandleLoginPageRoute).Methods("GET")
  r.HandleFunc("/login",controller.HandleLoginAction).Methods("POST")
  http.Handle("/",r)
  
  err := user.CreateUser("test", "test")
  if err != nil {
    log.Print(err)
  }
  err = user.VerifyUser("test", "tdest")
  if err != nil {
    log.Print(err)
  }
  
  // TODO: Use HTTPS
  http.ListenAndServe("127.0.0.1:80",nil)
}