package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"

	"github.com/metaclips/LetsTalk/backend/model"
	"github.com/metaclips/LetsTalk/backend/values"
)

var homepageTmpl, loginTmpl *template.Template

func init() {
	var terr error
	// Use (%%) instead of {{}} for templates.
	homepageTmpl, terr = template.New("home.html").Delims("(%", "%)").ParseFiles(
		"backend/views/homepage/home.html",
		"backend/views/homepage/components/SideBar.vue", "backend/views/homepage/components/ChattingComponent.vue", "backend/views/homepage/components/CallUI.vue")

	if terr != nil {
		log.Fatalln("error parsing homepage templates", terr)
	}

	loginTmpl, terr = template.New("login.html").Delims("(%", "%)").ParseFiles("backend/views/loginpage/login.html")
	if terr != nil {
		log.Fatalln("error parsing login templates", terr)
	}
}

// HomePage is a GET request that validates user credentials and load homepage templates.
func HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := validatUser(w, r)

	switch data.(type) {
	case error:
		log.Println("could not log user in", data)
		http.Redirect(w, r, "/login", 302)

	default:
		if err := homepageTmpl.Execute(w, data); err != nil {
			log.Println(err)
		}
	}

}

// HomePageLoginGet loads login page for users to login. Cookies are initially validated if user is already logged in.
func HomePageLoginGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := validatUser(w, r)

	// On validate user, if users has initially logged in execute homepage template else execute login template.
	switch data.(type) {
	case error:
		data := setLoginDetails(false, false, "", "/login")

		if err := loginTmpl.Execute(w, data); err != nil {
			log.Println(err)
		}

	default:
		if err := homepageTmpl.Execute(w, data); err != nil {
			log.Println(err)
		}
	}
}

func HomePageLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}

	email := strings.ToLower(r.FormValue("username"))
	password := r.FormValue("password")

	token, err := model.User{
		Email: email,
	}.CreateUserLogin(password, w)

	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	sendResponse(w, token)
}

func validatUser(w http.ResponseWriter, r *http.Request) interface{} {
	cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}
	if err := cookie.CheckCookie(r, w); err != nil {
		return err
	}

	data := struct {
		Email, UUID, Name string
	}{
		cookie.Email, cookie.Data.UUID,
		values.MapEmailToName[cookie.Email],
	}

	return data
}

func setLoginDetails(errors, isAdmin bool, errorDetail, link string) struct {
	SigninError, Admin bool
	Login, ErrorDetail string
} {

	return struct {
		SigninError, Admin bool
		Login, ErrorDetail string
	}{
		errors,
		isAdmin,
		link,
		errorDetail,
	}
}

func sendResponse(w http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Errorf("error sending response, err: %s \n", err.Error())
	}
}
