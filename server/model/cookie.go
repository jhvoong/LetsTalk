package model

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/metaclips/LetsTalk/values"
	"go.mongodb.org/mongo-driver/bson"
)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func (b CookieDetail) CreateCookie(w http.ResponseWriter) error {
	exitTime := time.Now().Add(time.Hour * 2).Local()
	b.Data.ExitTime = exitTime
	b.Data.UUID = uuid.New().String()

	_, err := db.Collection(b.Collection).UpdateOne(ctx, bson.M{"_id": b.Email},
		bson.M{"$set": bson.M{"loginUUID": b.Data.UUID, "expires": exitTime}})
	if err != nil {
		return err
	}

	encoded, err := cookieHandler.Encode(b.CookieName, b.Data)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     b.CookieName,
		Value:    encoded,
		Expires:  exitTime,
		SameSite: http.SameSiteStrictMode,
		Secure:   true, // Cookie is set to secure that is https so non-https would be dropped.
		Path:     b.Path,
	}

	http.SetCookie(w, cookie)
	return nil
}

func (b *CookieDetail) CheckCookie(r *http.Request, w http.ResponseWriter) error {
	cookie, err := r.Cookie(b.CookieName)
	if err != nil {
		return err
	}

	if err := cookieHandler.Decode(b.CookieName, cookie.Value, &b.Data); err != nil {
		// Reset cookies if cookie validation fails.
		http.SetCookie(w, &http.Cookie{
			Name:    b.CookieName,
			Expires: time.Now(),
		})

		return err
	}

	if b.Data.ExitTime.Before(time.Now().Local()) {
		return values.ErrCookieExpired
	}

	b.Email = b.Data.Email
	result := db.Collection(b.Collection).FindOne(ctx, bson.M{"_id": b.Email})

	cookieUUID := struct {
		LoginUUID string    `json:"loginUUID"`
		Expires   time.Time `json:"expires"`
	}{}

	if err = result.Decode(&cookieUUID); err != nil {
		return err
	}

	if cookieUUID.LoginUUID != b.Data.UUID {
		return values.ErrIncorrectUUID
	}

	if cookieUUID.Expires.Sub(b.Data.ExitTime).Seconds() > 0 {
		return values.ErrInvalidExpiryTime
	}

	return nil
}
