package main

import (
	. "booking-api/db"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

func createJwt(payload jwt.MapClaims) string {
	payload["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}

func main() {
	// wipe existing data
	DB.MustExec(`TRUNCATE TABLE "user" CASCADE`)

	// seed.exe with fresh data
	id1 := "3a39d390-ec5a-4795-b584-ca51ecd73c28"
	id2 := "48457cf0-7411-4b3f-87da-78ddcae82589"
	jwt1 := createJwt(jwt.MapClaims{
		"id": id1,
	})
	jwt2 := createJwt(jwt.MapClaims{
		"id": id2,
	})
	DB.MustExec(`
	INSERT INTO "user"(id, name) 
	VALUES 
	  ($1, 'bob'),
		($2, 'jim')
	RETURNING id
	`, id1, id2)

	// output jwts for generated users
	log.Println(jwt1)
	log.Println(jwt2)
}
