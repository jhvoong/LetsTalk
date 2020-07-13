package controller

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/metaclips/LetsTalk/model"
	"github.com/metaclips/LetsTalk/values"
	log "github.com/sirupsen/logrus"
)

var homepageTmpl, loginTmpl *template.Template

func init() {
	var terr error
	// Use (%%) instead of {{}} for templates.
	homepageTmpl, terr = template.New("home.html").Delims("(%", "%)").ParseFiles(
		"views/homepage/home.html",
		"views/homepage/components/SideBar.vue", "views/homepage/components/ChattingComponent.vue", "views/homepage/components/CallUI.vue")

	if terr != nil {
		log.Fatalln("error parsing homepage templates", terr)
	}

	loginTmpl, terr = template.New("login.html").Delims("(%", "%)").ParseFiles("views/loginpage/login.html")
	if terr != nil {
		log.Fatalln("error parsing login templates", terr)
	}
}

// HomePage is a GET request that validates user credentials and load homepage templates.
func HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := validatUser(w, r)

	switch data.(type) {
	case error:
		log.Errorln("could not log user in", data)
		http.Redirect(w, r, "/login", 302)

	default:
		if err := homepageTmpl.Execute(w, data); err != nil {
			log.Errorln("error executing homepage template on HomePage", err)
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
			log.Warningln("error executing login template on HompageLoginGet", err)
		}

	default:
		if err := homepageTmpl.Execute(w, data); err != nil {
			log.Warningln("error executing homepage template on HomepageLoginGet", err)
		}
	}
}

func HomePageLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	email := r.FormValue("email")
	email = strings.ToLower(email)
	password := r.FormValue("password")

	err := model.User{
		Email: email,
	}.CreateUserLogin(password, w)

	if err != nil {
		data := setLoginDetails(true, false, "Username or password invalid.", "/login")

		if err := loginTmpl.Execute(w, data); err != nil {
			log.Warningln("error executing login template on HomePageLoginPost", err)
		}
		return
	}

	http.Redirect(w, r, "/", 302)
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
