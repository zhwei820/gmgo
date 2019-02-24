package main

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/narup/gmgo"
	"log"
)

var TestDB gmgo.Db

//####################
type User struct {
	Id    bson.ObjectId `json:"_id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Email string        `json:"email" bson:"email"`
	Sex   string        `json:"sex" bson:"sex"`
}

// Each of your data model that needs to be persisted should implment gmgo.Document interface
func (user User) CollectionName() string {
	return "user"
}

//####################

func saveNewUser() {
	session := TestDB.Session()
	defer session.Close()

	user := &User{Name: "Puran", Email: "puran@xyz.com"}
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
func findAllUsers() {
	session := TestDB.Session()
	defer session.Close()

	users, err := session.FindAll(gmgo.Q{}, new(User)) //Note user pointer is passed to identify the collection type etc.
	if err != nil {
		fmt.Printf("Error fetching users %s", err)
	} else {
		for _, user := range users.([]*User) {
			println(user.Id.Hex() + " -- " + user.Email)
		}
	}
}

func findUsingIterator() ([]*User, error) {
	session := TestDB.Session()
	defer session.Close()

	users := make([]*User, 0)

	itr := session.DocumentIterator(gmgo.Q{"name": "Puran"}, "user")
	itr.Load(gmgo.IteratorConfig{Limit: 20, SortBy: []string{"-_id"}})

	result, err := itr.All(new(User))
	if err != nil {
		println(err)
	}
	users1 := result.([]*User)
	for _, user := range users1 {
		println(user.Id.Hex() + " -- " + user.Email)
	}

	return users, nil
}

func setupDB() {
	if err := gmgo.Setup(gmgo.DbConfig{
		HostURL: "mongodb://localhost:27017/userdb",
		DBName:  "userdb"}); err != nil {
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

	println("saveNewUser")
	for ii := 0; ii < 10000; ii++ {
		saveNewUser()
	}

	println("findUser")
	user := findUser("5c7250d7f7a26520feb67258")
	if user != nil {
		println(fmt.Sprintf("User name:%v", user.Name))
	} else {
		println("Couldnt find user")
	}

	println("findAllUsers")
	findAllUsers()

	println("findUsingIterator")
	findUsingIterator()
}
