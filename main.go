package main

import (
	"fmt"
	"log"

	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	api := iris.New()
	api.Use(&MongoDBMiddleware{})
	api.Get("/someGet", getting)
	api.Listen(":3000")
}

type Person struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Name  string
	Phone string
}

func getting(ctx *iris.Context) {
	dbconn, ok := ctx.Get("databaseConn").(*mgo.Session)
	if !ok {
		fmt.Println("GG")
	}
	defer dbconn.Close()

	err := dbconn.Ping()
	if err != nil {
		log.Println("Ping DB failed.")
	}
	c := dbconn.DB("test").C("people")

	// Query All
	var results []Person
	err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)

	ctx.Write("Hi %s\n", "iris")
}

type MongoDBMiddleware struct {
	// your 'stateless' fields here
}

func (m *MongoDBMiddleware) Serve(ctx *iris.Context) {
	// db init
	mongodb := "127.0.0.1"
	session, err := mgo.Dial(mongodb)
	if err != nil {
		log.Println("cannot connect to mongo, error:", err)
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	ctx.Set("databaseConn", session)
	ctx.Next()
}
