package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	//"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"math/rand"
	"net/http"
	"time"
)

type Data struct{
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Location string	`json:"location" bson:"location"`
	Birthday string `json:"birthday" bson:"birthday"`
}

var(
	client = mongoInit()
	ctx = context.Background()
	col = client.Database("testing").Collection("data")
)

func mongoInit() *mongo.Client{
	// Set client options
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://aldhirafe:pohTo5-wotmog-fidhyq@mongo-trial-data-yg8bk.gcp.mongodb.net/test"))
	if err != nil {
		fmt.Println("Can't Connect to MongoDB!")
		log.Fatal(err)
	}
	return client
}

func checkMongoConnect(){
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("MongoDB not Connected!")
		log.Fatal(err)
	}
	fmt.Println("Connected!")
}

func justAnAPI(){
	fmt.Println("API Started!")
	r := mux.NewRouter()

	r.HandleFunc("/data", createData).Methods("POST")
	r.HandleFunc("/data", getDatas).Methods("GET")
	r.HandleFunc("/data/{id}", getData).Methods("GET")
	r.HandleFunc("/data/{id}", deleteData).Methods("DELETE")
	r.HandleFunc("/data/{id}", updateData).Methods("PUT")

	http.ListenAndServe(":8080", r)
}

func createPerson(data Data) map[string]interface{}{

	res, err := col.InsertOne(ctx, data)

	response := make(map[string]interface{})

	if err == nil {
		response["message"] = "success"
		response["id"] = res.InsertedID
		response["status"] = true
	} else {
		response["message"] = "failed"
		response["id"] = nil
		response["status"] = false
	}

	return response
}

func getPeople() map[string]interface{} {
	response := make(map[string]interface{})
	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		response["message"] = "failed"
		response["data"] = nil
		response["status"] = false
		log.Fatal(err)
	} else {
		defer cur.Close(ctx)

		var getdata []Data

		for cur.Next(ctx){
			var data Data
			_ = cur.Decode(&data)
			getdata = append(getdata, data)
		}
		response["message"] = "success"
		response["data"] = getdata
		response["status"] = true
	}

	return response
}

func getPerson(id primitive.ObjectID) map[string]interface{}{
	var data Data
	response := make(map[string]interface{})
	err := col.FindOne(ctx, bson.M{"_id" : id}).Decode(&data)
	if err != nil {
		response["message"] = "failed"
		response["data"] = nil
		response["status"] = false
	} else {
		response["message"] = "success"
		response["data"] = data
		response["status"] = true
	}
	return response
}

func deletePerson(id primitive.ObjectID) map[string]interface{}{
	response := make(map[string]interface{})
	res, err := col.DeleteOne(ctx, bson.M{"_id" : id})
	if err != nil {
		response["message"] = "failed"
		response["data"] = nil
		response["status"] = false
	} else {
		response["message"] = "deleted"
		response["data"] = res.DeletedCount
		response["status"] = true
	}
	return response
}

func updatePerson(data Data, id primitive.ObjectID) map[string]interface{}{
	response := make(map[string]interface{})
	data.ID = id
	_ , err := col.UpdateOne(ctx, bson.M{"_id": bson.M{"$eq": id}}, bson.M{"$set": data})
	if err != nil {
		response["message"] = "failed"
		response["data"] = nil
		response["status"] = false
		fmt.Print(err)
	} else {
		response["message"] = "updated"
		response["data"] = data
		response["status"] = true
	}

	return response
}

func createData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data Data
	_ = json.NewDecoder(r.Body).Decode(&data)

	json.NewEncoder(w).Encode(createPerson(data))
}

func getDatas(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getPeople())
}

func getData(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	json.NewEncoder(w).Encode(getPerson(id))
}

func deleteData (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	json.NewEncoder(w).Encode(deletePerson(id))
}

func updateData (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var data Data
	_ = json.NewDecoder(r.Body).Decode(&data)
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	json.NewEncoder(w).Encode(updatePerson(data, id))
}


func main(){
	rand.Seed(time.Now().UTC().UnixNano())
	checkMongoConnect()
	//API Connection
	justAnAPI()
}