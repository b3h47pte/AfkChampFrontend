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
  "github.com/gorilla/mux"
  "log"
)

type AdminUserTemplateData struct {
  Data controller.BaseTemplateData
  Users []user.UserEntry
  SelectedUser user.UserEntry
  IsNewUser bool
  
  UserNameCharLimit int
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

// 'HandleAdminEditUserRoute' takes in a user and presents properties that we can modify.
func HandleAdminEditUserRoute(w http.ResponseWriter, r *http.Request) {
  handleAdminUserNewEditPage(w, r, false)
}

// 'handleAdminUserNewEditPage' prepares the new/edit page for this specific case.
// Note: I don't actually have the functionality here to create a new user...that should just be through the register page.
func handleAdminUserNewEditPage(w http.ResponseWriter, r *http.Request, isNewUserRequest bool) {
  if err := RequireAdminRelogin(w, r); err != nil {
    return
  }
  
  // Retrieve relevant user
  userVars := mux.Vars(r)
  userid, err := strconv.ParseInt(userVars["userid"], 10, 64)
  if err != nil {
    log.Print(err)
    http.Redirect(w, r, "/admin/user", http.StatusFound)
    return
  }
  
  selectedUser, err := user.RetrieveUser(userid)
  if err != nil {
    log.Print(err)
    http.Redirect(w, r, "/admin/user", http.StatusFound)
    return
  }
  
  t := CreateBaseUserAdminTemplateData()
  t.SelectedUser = *selectedUser
  controller.TemplateMapping["admin/user/newedit.html"].ExecuteTemplate(w, "tbase", t)
}

// 'CreateBaseUserAdminTemplateData' creates the template data for rendering.
func CreateBaseUserAdminTemplateData() *AdminUserTemplateData {
  t := AdminUserTemplateData{Data: controller.CreateTemplateData(),
                            UserNameCharLimit: user.MaxUsernameLength}
  return &t
}