package main

import (
	"log"
	"time"

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
	Date  int64         `bson:"date"`
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
	err = c.Insert(&Post{Id: i, Image: img, Date: time.Now().Unix()})
	if err != nil {
		return "", err
	}

	log.Println("Successfully inserted", i.Hex())

	return i, nil
}

func updateCode(id, code string) error {
	s, err := open()
	if err != nil {
		return err
	}
	defer s.Close()

	c := s.DB(DBNAME).C(COLLECTION)

	// Get existing post and change it
	post, err := getPost(id)
	if err != nil {
		return err
	}
	post.Code = code

	// Push that update
	_, err = c.UpsertId(bson.ObjectIdHex(id), post)
	if err != nil {
		return err
	}

	log.Println("Successfully updated code for", id)

	return nil
}

func getPost(id string) (*Post, error) {
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

	return &result, nil
}

func getImage(id string) ([]byte, error) {
	post, err := getPost(id)
	if err != nil {
		return nil, err
	}
	return post.Image, err
}
