package main

import (
	"AfkChampFrontend/controller"
	"AfkChampFrontend/controller/about"
	"AfkChampFrontend/controller/admin"
	"AfkChampFrontend/controller/live"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type AfkChampHandler func(w http.ResponseWriter, req *http.Request)

func main() {
	// Setup logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Setup the routing for the frontend.
	r := mux.NewRouter()
	// Dynamic Content
	r.HandleFunc("/", controller.HandleHomeRoute).Methods("GET")

	// ADMIN
	r.HandleFunc("/admin", admin.HandleAdminRoute).Methods("GET")
	r.HandleFunc("/admin/game", admin.HandleAdminGamePageRoute).Methods("GET")
	r.HandleFunc("/admin/game", admin.HandleNewEditGamePost).Methods("POST")
	r.HandleFunc("/admin/game/new", admin.HandleAdminGameNewRoute).Methods("GET")
	r.HandleFunc("/admin/game/{gameName}", admin.HandleAdminGameEditRoute).Methods("GET")
	r.HandleFunc("/admin/game/{gameName}/delete", admin.HandleAdminGameDeleteRoute).Methods("GET")
	r.HandleFunc("/admin/user", admin.HandleAdminUserIndexPageRoute).Methods("GET")
	r.HandleFunc("/admin/user/{userid}", admin.HandleAdminEditUserRoute).Methods("GET")
	r.HandleFunc("/admin/user/{userid}/delete", admin.HandleAdminDeleteUsePage).Methods("GET")
	r.HandleFunc("/admin/user", admin.HandleAdminUserNewEditPost).Methods("POST")
	r.HandleFunc("/admin/event", admin.HandleAdminEventNewEditPost).Methods("POST")
	r.HandleFunc("/admin/event", admin.HandleAdminEventIndexRoute).Methods("GET")
	r.HandleFunc("/admin/event/new", admin.HandleAdminEventNewRoute).Methods("GET")
	r.HandleFunc("/admin/event/{eventShorthand}", admin.HandleAdminEventEditRoute).Methods("GET")
	r.HandleFunc("/admin/event/{eventShorthand}/delete", admin.HandleAdminEventDeleteRoute).Methods("GET")

	// LIVE STATS
	r.HandleFunc("/live/{eventShorthand}", live.HandleLiveStatsPageRoute).Methods("GET")

	// MAIN PAGE
	r.HandleFunc("/login", controller.HandleLoginPageRoute).Methods("GET")
	r.HandleFunc("/login", controller.HandleLoginAction).Methods("POST")
	r.HandleFunc("/register", controller.HandleRegisterPageRoute).Methods("GET")
	r.HandleFunc("/register", controller.HandleRegisterAction).Methods("POST")
	r.HandleFunc("/logout", controller.HandleLogoutPageRoute).Methods("GET")
	r.HandleFunc("/about", about.HandleAboutRoute).Methods("GET")

	// Static Content
	r.PathPrefix("/javascript/").Handler(http.StripPrefix("/javascript/", http.FileServer(http.Dir("./javascript/"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))

	// Error Pages
	r.NotFoundHandler = http.HandlerFunc(controller.Handle404Page)

	http.Handle("/", r)

	// TODO: Use HTTPS
	// TODO: Listen on the right IP address (config?)
	http.ListenAndServe("127.0.0.1:80", nil)
}
