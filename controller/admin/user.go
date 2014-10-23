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
  "AfkChampFrontend/utility"
)

type AdminUserTemplateData struct {
  Data controller.BaseTemplateData
  Users []user.UserEntry
  SelectedUser user.UserEntry
  IsNewUser bool
  
  UserNameCharLimit int
}

type AdminUserErrorCode int
const (
  errorUserNoError AdminUserErrorCode = iota
  errorUserInvalidOperation
  errorUserUnspecifiedError
)

type AdminUserPostData struct {
  IsNew bool
  User user.UserEntry
}

type AdminUserResponseData struct {
  ErrorCode AdminUserErrorCode 
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

// 'HandleAdminUserNewEditPost' handles POST requests for new/edit user requests from the admin.
func HandleAdminUserNewEditPost(w http.ResponseWriter, r *http.Request) {
  if err := RequireAdminRelogin(w, r); err != nil {
    return
  }
  
  userData := AdminUserPostData{}
  err := utility.ReadJsonFromRequestBodyStruct(r, &userData)
  if err != nil {
    log.Print(err)
    AdminUserRespondJsonError(errorUserUnspecifiedError, w)
    return
  }
  
  // "New" requests are bogus. Don't let them happen. [FOR NOW]
  if userData.IsNew {
    AdminUserRespondJsonError(errorUserInvalidOperation, w)
    return
  }
  
  // Modify the user as appropriate.
  if err = user.UpdateUser(userData.User.UserId, &userData.User); err != nil {
    log.Print(err)
    AdminUserRespondJsonError(errorUserUnspecifiedError, w)
    return 
  }
  
  AdminUserRespondJsonError(errorUserNoError, w)
}

// 'CreateBaseUserAdminTemplateData' creates the template data for rendering.
func CreateBaseUserAdminTemplateData() *AdminUserTemplateData {
  t := AdminUserTemplateData{Data: controller.CreateTemplateData(),
                            UserNameCharLimit: user.MaxUsernameLength}
  return &t
}

// AdminUserRespondJsonError takes in an error code and passes it back to the client in the form of a JSON response.
func AdminUserRespondJsonError(errorCode AdminUserErrorCode, w http.ResponseWriter) {
  response := AdminUserResponseData{ErrorCode: errorCode}
  if errorCode != errorUserNoError {
    htmlErrCode := getErrorCodeFromUserError(errorCode)
    http.Error(w, "", htmlErrCode)
  }
  
  // If any error happens here, then the only thing we can redirect the user to is an error page.
  err := utility.WriteJsonToResponse(w, response)
  if err != nil {
    log.Print(err)
    return
  }
}

// 'getErrorCodeFromUserError' takes in a game error code and returns a HTTP error code along with it.
func getErrorCodeFromUserError(errorCode AdminUserErrorCode) int {
  switch errorCode {
  case errorUserUnspecifiedError, errorUserInvalidOperation:
    return 500
  default:
    return 200
  }
  return 200
}