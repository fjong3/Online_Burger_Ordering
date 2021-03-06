package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/unrolled/render"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/satori/go.uuid"
	"net"
	"os"
	"strings"
)

var mongodb_server = os.Getenv("MONGO_SERVER")
var mongodb_database = os.Getenv("MONGO_DATABASE")
var mongodb_collection = os.Getenv("MONGO_COLLECTION")
var mongo_admin_database = os.Getenv("MONGO_ADMIN_DATABASE")
var mongo_username = os.Getenv("MONGO_USERNAME")
var mongo_password = os.Getenv("MONGO_PASS")

func MenuServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	router := mux.NewRouter()
	initRoutes(router, formatter)
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD","DELETE", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	n.UseHandler(handlers.CORS(allowedHeaders,allowedMethods , allowedOrigins)(router))
	return n
}

func initRoutes(router *mux.Router, formatter *render.Render) {
	router.HandleFunc("/users", GetAllUser).Methods("GET")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/users/signup", CreateUser).Methods("POST")
	router.HandleFunc("/users/signin", UserSignIn).Methods("POST")
	router.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	router.HandleFunc("/users/test/ping", checkPing(formatter)).Methods("GET")
}

func checkPing(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		message := "Burger Users API Server Working on machine: " + getSystemIp()
		formatter.JSON(w, http.StatusOK, struct{ Test string }{message})
	}
}

func getSystemIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
		return "" 
	}
    defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).String()
	address := strings.Split(localAddr, ":")
    fmt.Println("address: ", address[0])
    return address[0]
}

func GetUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	var m User
    _ = json.NewDecoder(req.Body).Decode(&m)		
    fmt.Println("Get data of user: ", params["id"])
	session, err := mgo.Dial(mongodb_server)
    if err != nil {
		message := struct {Message string}{"Some error occured while connecting to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	err = session.DB(mongo_admin_database).Login(mongo_username, mongo_password)
	if err != nil {
		message := struct {Message string}{"Some error occured while login to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    c := session.DB(mongodb_database).C(mongodb_collection)
    query := bson.M{"id" : params["id"]}
    var result bson.M
    err = c.Find(query).One(&result)
    if err != nil && err != mgo.ErrNotFound{
		message := struct {Message string}{"Some error occured while querying to database!!"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(message)
		return
    }else if err == mgo.ErrNotFound{
		message := struct {Message string}{"User not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(message)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func GetAllUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	session, err := mgo.Dial(mongodb_server)
    if err != nil {
		message := struct {Message string}{"Some error occured while connecting to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	err = session.DB(mongo_admin_database).Login(mongo_username, mongo_password)
	if err != nil {
		message := struct {Message string}{"Some error occured while login to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    c := session.DB(mongodb_database).C(mongodb_collection)
    query := bson.M{}
    var result []bson.M
    err = c.Find(query).All(&result)
    if err != nil {
		message := struct {Message string}{"No users were found!!"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(message)
		return
    }
	json.NewEncoder(w).Encode(result)
}

func CreateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person User
	_ = json.NewDecoder(req.Body).Decode(&person)
	unqueId := uuid.Must(uuid.NewV4())
	person.Id = unqueId.String()
	session, err := mgo.Dial(mongodb_server)
    if err != nil {
		message := struct {Message string}{"Some error occured while connecting to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	err = session.DB(mongo_admin_database).Login(mongo_username, mongo_password)
	if err != nil {
		message := struct {Message string}{"Some error occured while login to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"email" : person.Email}
    var result bson.M
	err = c.Find(query).One(&result)
	if err != nil && err != mgo.ErrNotFound{
		message := struct {Message string}{"Some error occured while querying to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}else if result != nil {
		message := struct {Message string}{"User already exists!!"}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(message)
		return
	}
    err = c.Insert(person)
    if err != nil {
		message := struct {Message string}{"Some error occured while querying to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	session, err := mgo.Dial(mongodb_server)
    if err != nil {
		message := struct {Message string}{"Some error occured while connecting to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	err = session.DB(mongo_admin_database).Login(mongo_username, mongo_password)
	if err != nil {
		message := struct {Message string}{"Some error occured while login to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"id":params["id"]}
    err = c.Remove(query)
    if err != nil && err != mgo.ErrNotFound {
		message := struct {Message string}{"Some error occured while querying to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		fmt.Println("error:"+ err.Error())
		return
    }else if err == mgo.ErrNotFound {
		message := struct {Message string}{"user not found!!"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(message)
		return
	}
	json.NewEncoder(w).Encode(struct {Message string }{"user with id:"+params["id"]+" was deleted"})
}
func UpdateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person User
	_ = json.NewDecoder(req.Body).Decode(&person)
	params := mux.Vars(req)
	session, err := mgo.Dial(mongodb_server)
    if err != nil {
		message := struct {Message string}{"Some error occured while connecting to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	err = session.DB(mongo_admin_database).Login(mongo_username, mongo_password)
	if err != nil {
		message := struct {Message string}{"Some error occured while login to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"id":params["id"]}
	updator := bson.M{
				"$set": bson.M{
						"firstname": person.Firstname,
						"lastname":person.Lastname,
						"address":person.Address,
						"password":person.Password}}
    err = c.Update(query, updator)
    if err != nil && err != mgo.ErrNotFound{
		message := struct {Message string}{"Some error occured while querying to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
    }else if err == mgo.ErrNotFound{
		message := struct {Message string}{"User not found!!"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(message)
		return
	}
	json.NewEncoder(w).Encode(struct {Message string }{"user with id:"+params["id"]+" was Updated"})
}
func UserSignIn(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person User
	_ = json.NewDecoder(req.Body).Decode(&person)
	session, err := mgo.Dial(mongodb_server)
    if err != nil {
		message := struct {Message string}{"Some error occured while connecting to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	err = session.DB(mongo_admin_database).Login(mongo_username, mongo_password)
	if err != nil {
		message := struct {Message string}{"Some error occured while login to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongodb_database).C(mongodb_collection)
	query := bson.M{"email":person.Email,
					"password":person.Password}
	var result User
    err = c.Find(query).One(&result)
    if err != nil && err != mgo.ErrNotFound{
		message := struct {Message string}{"Some error occured while querying to database!!"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(message)
		return
	}
	if err == mgo.ErrNotFound {
		message := struct {Message string}{"Login Failed"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(message)
		return
	}
	userData := bson.M{ 
						"email":result.Email,
						"firstName":result.Firstname,
						"lastName":result.Lastname,
						"address":result.Address,
						"id":result.Id}
	json.NewEncoder(w).Encode(userData)
}
