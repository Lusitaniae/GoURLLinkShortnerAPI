package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var (
	dynamoRegion   = flag.String("dynamo-region", "eu-west-1", "AWS region where DynamoDB database is hosted.")
	dynamoEndpoint = flag.String("dynamo-endpoint", "http://127.0.0.1:8000", "DynamoDB database address.")
	dynamoTable    = flag.String("dynamo-table", "url_shortener", "DynamoDB Table to store shortened urls.")
	redisEndpoint  = flag.String("redis-endpoint", "127.0.0.1:6379", "Redis database address.")
)

/*
	This is the entry point for the API, the purpose of the API is to shorten url links via a REST HTTP request
	Once the code is built and running, to create a short url to a long url mapping, send a POST request to http://localhost:5100/Create, the POST request should include the shortUrl and the longUrl as follows:
	{'shorturl':'cosmosfading','longurl':'https://cosmosmagazine.com/space/universe-slowly-fading-away'}

	You could then consume a shorturl by issuing a GET request to http://localhost:5100/<the short url>
	The program makes use of the gorilla mux library for routing as well as the mgo library to interface with mongo database
*/
func main() {

	flag.Parse()
	//Create a new API shortner API
	LinkShortener := NewUrlLinkShortenerAPI()
	//Create the needed routes for the API
	routes := CreateRoutes(LinkShortener)
	//Initiate the API routers
	router := NewLinkShortenerRouter(routes)
	//This will start the web server on local port 5100
	srv := &http.Server{
		Handler: router,
		Addr:    ":5100",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
