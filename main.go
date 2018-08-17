// Unity Test API Version 2.0.2
// Developed by Alireza Abdelahi @ 2018
package main

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/nats-io/go-nats"
	"github.com/valyala/fasthttp"
	gojsonschema "github.com/xeipuuv/gojsonschema"
)

//Schema validation data (draft 7)
//In real production and for better code style is better to store it in uri and cache it to reduce network latancy while it fetch from central location
var payloadSchema string

type server struct {
	nc *nats.Conn
}

var version = "2.0.2"

// main function to boot up everything
func main() {

	var s server
	var err error
	uri := os.Getenv("NATS_URI")
	//uri := "nats://demo.nats.io:4222"
	fmt.Println("NAT_URI=" + uri)
	buf, err := ioutil.ReadFile("payloadSchema.json")
	if err == nil {
		payloadSchema = string(buf)
	} else {
		log.Fatal("Error Loading file", err)
	}
	//try to connect to NATS server (retry 5 times)
	for i := 0; i < 5; i++ {
		nc, err := nats.Connect(uri)
		if err == nil {
			s.nc = nc
			break
		}
		//fmt.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}
	//fmt.Println("Connected to NATS at:", s.nc.ConnectedUrl())

	//Define http routings
	router := fasthttprouter.New()
	router.GET("/", indexHandler)
	router.GET("/healthcheck", healthzHandler)
	router.GET("/loaderio-fd520a492327c3a48606330949d7e368/", verify)
	router.POST("/api/v2/unitytestapi", s.registerMessage)

	//Run Http listener
	fmt.Printf("Server version 2.0.1 is listening on port 80...")

	log.Fatal(fasthttp.ListenAndServe(":80", router.Handler))
}

// request handler in fasthttp style, i.e. just plain function.
func (s server) registerMessage(ctx *fasthttp.RequestCtx) {
	//Validate json input by defined schema
	schemaLoader := gojsonschema.NewStringLoader(payloadSchema)
	documentLoader := gojsonschema.NewStringLoader(string(ctx.PostBody()))
	jsonValidationResult, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if jsonValidationResult.Valid() {
		fmt.Fprintln(ctx, "The json parameter is valid. ")
		if s.nc.IsConnected() {
			fmt.Fprintln(ctx, "Connected to "+s.nc.ConnectedUrl())
			// Simple Publish
			puberror := s.nc.Publish("messageTopic", ctx.PostBody())
			if puberror == nil {
				fmt.Fprintln(ctx, "message has been Published")
			} else {
				fmt.Fprintln(ctx, "message has not been Published")
			}
		}
	} else {
		fmt.Fprint(ctx, "The json parameter is not valid. see errors :\n")
		//for _, desc := range result.Errors() {
		//	fmt.Printf("- %s\n", desc)
		//}
	}

}

// Display API version
func indexHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Unity Validation API Test version "+version)
}

//health check endpoint to ensure service is up
func healthzHandler(ctx *fasthttp.RequestCtx) {
	fmt.Println("HealthCheck request recived")
	fmt.Fprintln(ctx, "OK")
}

//verify check endpoint to load.io verify
func verify(ctx *fasthttp.RequestCtx) {
	fmt.Fprintln(ctx, "loaderio-fd520a492327c3a48606330949d7e368")
}
