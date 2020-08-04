package controller

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/metaclips/LetsTalk/model"
	"github.com/metaclips/LetsTalk/values"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func AdminLoginPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	data := setLoginDetails(false, true, "", "/admin/login/")

	admin := model.Admin{StaffDetails: model.User{Email: email}}
	if err := admin.CheckAdminDetails(password); err != nil {
		data.SigninError = true
		data.ErrorDetail = values.ErrInvalidDetails.Error()

		if err := loginTmpl.Execute(w, data); err != nil {
			log.Warningln(err)
		}

		return
	}

	cookie := model.CookieDetail{
		Email:      admin.StaffDetails.Email,
		Collection: values.AdminCollectionName,
		CookieName: values.AdminCookieName,
		Path:       "/admin",
		Data: model.CookieData{
			Super: admin.Super,
			Email: admin.StaffDetails.Email,
		},
	}

	if err := cookie.CreateCookie(w); err != nil {
		log.Warningln("error creating cookies on admin login post", err)
		data.SigninError = true
		data.ErrorDetail = "server error"
		loginTmpl.Execute(w, data)
		return
	}

	http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
}

func AdminLoginGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.AdminCookieName, Collection: values.AdminCollectionName}
	if err := cookie.CheckCookie(r, w); err == nil {
		http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
		return
	}

	data := setLoginDetails(false, true, "", "/admin/login/")
	if err := loginTmpl.Execute(w, data); err != nil {
		log.Warningln("error executing login template", err)
	}
}

func AdminPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.AdminCookieName, Collection: values.AdminCollectionName}
	if err := cookie.CheckCookie(r, w); err != nil {
		log.Warningln("error checking admin cookies on AdminPage", err)
		http.Redirect(w, r, "/admin/login/", http.StatusBadRequest)
		return
	}

	data := struct {
		UploadSuccess bool
		Error         bool
	}{
		false,
		false,
	}

	tmpl, terr := template.New("admin.html").Delims("(%", "%)").ParseFiles("views/admin/admin.html", "views/admin/components/tabs.vue",
		"views/admin/components/adduser.vue", "views/admin/components/block.vue", "views/admin/components/messagescan.vue")

	if terr != nil {
		log.Warningln("could not load template in AdminPage function", terr)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Warningln("could not execute template in AdminPage function", err)
	}
}

func UploadUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.AdminCookieName, Collection: values.AdminCollectionName}
	err := cookie.CheckCookie(r, w)
	if err != nil {
		http.Redirect(w, r, "/adnin/login", http.StatusBadRequest)
		return
	}

	r.ParseForm()

	user := model.User{
		Email:        strings.ToLower(r.FormValue("email")),
		Name:         r.FormValue("name"),
		DOB:          r.FormValue("DOB"),
		Class:        r.FormValue("usersClass"),
		Faculty:      r.FormValue("faculty"),
		ParentEmail:  r.FormValue("parentEmail"),
		ParentNumber: r.FormValue("parentNumber"),
	}

	data := struct {
		UploadSuccess bool
		Error         bool
	}{
		true,
		false,
	}

	if err := user.UploadUser(); err != nil {
		data.UploadSuccess = false
		data.Error = true
	}

	tmpl, terr := template.New("admin.html").Delims("(%", "%)").ParseFiles("views/admin/admin.html", "views/admin/components/tabs.vue",
		"views/admin/components/adduser.vue", "views/admin/components/block.vue", "views/admin/components/messagescan.vue")
	if terr != nil {
		log.Warningln("could not load template in UploadUser function", terr)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Warningln("could not execute template in UploadUser function", err)
	}
}
