package admin

/*
 * user.go in the 'admin' package allows us to modify user properties. This should be different from the regular user's profile editing tool since there should be more
 * functionality here.
 */
 import(
  "net/http"
  "AfkChampFrontend/controller"
  "AfkChampFrontend/model/user"
   "strconv"
)

type AdminUserTemplateData struct {
  Data controller.BaseTemplateData
  Users []user.UserEntry
}


// 'HandleAdminUserIndexPageRoute' presents a list of all users and allows us to view/modify any user properties as necessary.
func HandleAdminUserIndexPageRoute(w http.ResponseWriter, r *http.Request) {
  if err := RequireAdminRelogin(w, r); err != nil {
    return
  }
  
  // Show users in pages. So figure out which page we want
  pageIdx, err := strconv.Atoi(r.FormValue("p"))
  if err != nil {
    pageIdx = 0
  }
  
  // Get a portion of the users. 
  allUsers, err := user.IndexUsers(pageIdx * DefaultPageSize, DefaultPageSize)
  if err != nil {
    allUsers = make([]user.UserEntry, 0, 0)
  }
  
  t := CreateBaseUserAdminTemplateData()
  t.Users = allUsers
  controller.TemplateMapping["admin/user/index.html"].ExecuteTemplate(w, "tbase", t)
}

func HandleAdminUserRoute(w http.ResponseWriter, r *http.Request) {
//  userVars := mux.Vars(r)
//  username, ok := userVars["username"]
}

// 'CreateBaseUserAdminTemplateData' creates the template data for rendering.
func CreateBaseUserAdminTemplateData() *AdminUserTemplateData {
  t := AdminUserTemplateData{Data: controller.CreateTemplateData()}
  return &t
}