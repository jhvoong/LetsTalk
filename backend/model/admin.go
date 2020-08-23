package model

import (
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/metaclips/LetsTalk/backend/values"
)

func (b *Admin) CheckAdminDetails(password string) error {
	result := db.Collection(values.AdminCollectionName).FindOne(ctx, bson.M{
		"_id": b.StaffDetails.Email,
	})

	if err := result.Decode(&b); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(b.StaffDetails.Password, []byte(password)); err != nil {
		return err
	}
	return nil
}

func (b *Admin) CreateAdminLogin(w http.ResponseWriter) (string, error) {
	cookie := CookieDetail{
		Email:      b.StaffDetails.Email,
		Collection: values.AdminCollectionName,
		CookieName: values.AdminCookieName,
		Path:       "/admin",

		Data: CookieData{
			Super: b.Super,
			Email: b.StaffDetails.Email,
		},
	}

	return cookie.GenerateCookie(w)
}

func (b *Admin) CreateAdmin() error {
	_, err := db.Collection(values.AdminCollectionName).InsertOne(ctx, b)
	return err
}

func (b User) CreateUser(r *http.Request) error {
	b.Name = strings.Title(b.Name)

	if names := strings.Split(b.Name, " "); len(names) > 1 && b.PasswordInString == "" {
		b.PasswordInString = names[0]
	}

	var err error
	b.Password, err = bcrypt.GenerateFromPassword([]byte(b.PasswordInString), values.DefaultCost)
	if err != nil {
		return err
	}

	values.MapEmailToName[b.Email] = b.Name
	_, err = db.Collection(values.UsersCollectionName).InsertOne(ctx, b)
	return err
}
