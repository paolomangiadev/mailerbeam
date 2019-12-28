package config

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

// Session database type
type Session struct {
	*mgo.Session
}

// NewDB opens a new connection to database
func NewDB(url string) (*Session, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	return &Session{session}, err
}
