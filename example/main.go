package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/narup/gmgo"
	"log"
	"time"
)
import "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var Dbconfig = gmgo.DbConfig{
	HostURL: "mongodb://localhost:27017/userdb",
	DBName:  "userdb"}

var TestDB gmgo.Db

//####################
type User struct {
	Id    bson.ObjectId `json:"_id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Email string        `json:"email" bson:"email"`
	Sex   string        `json:"sex" bson:"sex"`
	Birth time.Time     `json:"birth" bson:"birth"`
	Addr  Addr          `json:"addr" bson:"addr"`
	Car   []Car         `json:"car" bson:"car"`
}

type Addr struct {
	Home string `json:"home" bson:"home"`
	Work string `json:"work" bson:"work"`
}

type Car struct {
	Brand string `json:"brand" bson:"brand"`
	Type  string `json:"type" bson:"type"`
}

// Each of your data model that needs to be persisted should implment gmgo.Document interface
func (user User) CollectionName() string {
	return "user"
}

//####################

func EnsureIndex() {
	session := TestDB.Session()
	defer session.Close()

	c := session.Session.DB(Dbconfig.DBName).C(User{}.CollectionName())
	index := mgo.Index{
		Key:         []string{"-birth"},
		Unique:      false,
		Background:  true,
		Sparse:      false,
		ExpireAfter: 0,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	idx, _ := c.Indexes()
	for id, item := range idx {
		println(fmt.Sprintf("%v, %v", id, item))
	}
}

func saveNewUser() {
	session := TestDB.Session()
	defer session.Close()

	user := &User{
		Name:  "Puran",
		Email: "puran@xyz.com",
		Birth: time.Now(),
		Car:   []Car{{Brand: "dfdsf", Type: "test"}},
	}
	user.Id = bson.NewObjectId()
	err := session.Save(user)
	if err != nil {
		log.Fatalf("Error saving user : %s.\n", err)
	}

	println(fmt.Sprintf("User id %s", user.Id))
}

func findUser(userId string) *User {
	session := TestDB.Session()
	defer session.Close()

	user := new(User)
	if err := session.FindByID(userId, user); err != nil {
		return nil
	}
	println(user.Id.Hex() + " -- " + user.Email)
	return user
}

//Find all users
func findAllUsers() ([]*User, error) {
	session := TestDB.Session()
	defer session.Close()

	users, err := session.FindAll(gmgo.Q{}, new(User)) //Note user pointer is passed to identify the collection type etc.
	return users.([]*User), err
}

func findUsingIterator() ([]*User, error) {
	session := TestDB.Session()
	defer session.Close()

	itr := session.DocumentIterator(gmgo.Q{"name": "Puran"}, "user")
	itr.Load(gmgo.IteratorConfig{Limit: 20, SortBy: []string{"-_id"}})

	result, err := itr.All(new(User))

	return result.([]*User), err
}

func setupDB() {
	if err := gmgo.Setup(Dbconfig); err != nil {
		log.Fatalf("Database connection error : %s.\n", err)
		return
	}

	newDb, err := gmgo.Get("userdb")
	if err != nil {
		log.Fatalf("Db connection error : %s.\n", err)
	}
	TestDB = newDb

}

func main() {
	//setup Mongo database connection. You can setup multiple db connections
	setupDB()
	EnsureIndex()

	println("saveNewUser")
	//for ii := 0; ii < 10000; ii++ {
	saveNewUser()
	//}

	println("findUser")
	user := findUser("5c725e16f7a2652e03d77cd8")
	if user != nil {
		println(fmt.Sprintf("User: %v", user.Birth))
	} else {
		println("Couldnt find user")
	}

	data, _ := json.Marshal(user)
	println(string(data))
	json.Unmarshal(data, &user)
	println(fmt.Sprintf("User: %v", user.Birth))

	//
	println("findAllUsers")
	findAllUsers()

	println("findUsingIterator")
	findUsingIterator()
}
