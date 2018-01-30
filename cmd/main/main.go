//Package main specifies that this is an executable command in Go.
//Files under this package are executables to start order microservice and initialize postgres DB.
package main

import (
	"github.com/gin-gonic/gin"
	"order/app"
	"os"
	"io"
	"log"
)

//setUpRouter creates a default gin router with appropriate handlers for multiple REST API endpoints for Order service.
//The API version can be incremented to cater to a input or output change or significant change, and backward compatibility can be maintained
//for clients still using the older ver. of APIs. Currently, v1 ver. of APIs is created.
func setUpRouter() *gin.Engine {

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	v1 := router.Group("/v1/order")
	{
		v1.HEAD("/", app.Ping)
		v1.POST("/", app.CreateOrder)
		v1.GET("/:id", app.FetchOrder)
		v1.PUT("/:id", app.UpdateOrder)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	return router
}

//setUpLogger configures the gin.log file where all the HTTP requests and application
//logs will be written.
func setUpLogger() {

	// Logging to a file.
	f, _ := os.Create("../../log/gin.log")

	//Write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.SetOutput(gin.DefaultWriter)
	log.Println("Logger is setup for the microservice")
}

//main is called when the executable runs.
//Sets up a web server on port 8080.
func main() {
	setUpLogger()
	router := setUpRouter()
	router.Run()

	//A custom HTTP configuration can also be provided to start the server as shown below.
	/*
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
	*/
}