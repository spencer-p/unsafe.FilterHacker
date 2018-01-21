package main

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DBURL      = "localhost:27017"
	DBNAME     = "test"
	COLLECTION = "posts"
)

type Post struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Image []byte        `bson:"image"`
	Code  string        `bson:"code"`
	Date  int           `bson:"date"`
}

func open() (*mgo.Session, error) {
	session, err := mgo.Dial(DBURL)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func uploadImage(img []byte) (bson.ObjectId, error) {
	s, err := open()
	if err != nil {
		return "", err
	}
	defer s.Close()

	// Make a new id
	i := bson.NewObjectId()

	// Insert the image into the posts collection
	c := s.DB(DBNAME).C(COLLECTION)
	err = c.Insert(&Post{Id: i, Image: img})
	if err != nil {
		return "", err
	}

	log.Print("Successfully inserted ", i.Hex())

	return i, nil
}

func getImage(id string) ([]byte, error) {
	s, err := open()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	c := s.DB(DBNAME).C(COLLECTION)

	var result Post
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return nil, err
	}

	return result.Image, nil
}
