package controller

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/metaclips/LetsTalk/backend/model"
	"github.com/metaclips/LetsTalk/backend/values"
)

// ServeWs handles websocket requests from the peer, ensuring user is registered.
func ServeWs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}
	if err := cookie.CheckCookie(r, w); err != nil {
		http.Error(w, values.ErrAuthentication.Error(), 404)
		return
	}

	ws, err := model.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &model.Connection{Send: make(chan []byte, 256), WS: ws}

	s := model.Subscription{Conn: c, User: cookie.Email}
	model.HubConstruct.Register <- s
	log.Println("User", cookie.Email, "Connected")
	go s.ReadPump(cookie.Email)
	s.WritePump()
}
