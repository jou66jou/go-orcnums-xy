package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/otiai10/marmoset"

	"github.com/jou66jou/go-orcnums-xy/controllers"
	"github.com/jou66jou/go-orcnums-xy/filters"
)

var logger *log.Logger
var port = flag.String("p", "8080", "server port")

func main() {
	flag.Parse()
	marmoset.LoadViews("./app/views")

	r := marmoset.NewRouter()
	// API
	r.GET("/status", controllers.Status)
	r.POST("/base64", controllers.Base64)
	r.POST("/file", controllers.FileUpload)
	// Sample Page
	r.GET("/", controllers.Index)
	r.Static("/assets", "./app/assets")

	logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", "ocrserver"), 0)
	r.Apply(&filters.LogFilter{Logger: logger})

	// port := os.Getenv("PORT")
	if port == nil {
		logger.Fatalln("Required env `PORT` is not specified.")
	}
	logger.Printf("listening on port %s", *port)
	if err := http.ListenAndServe(":"+*port, r); err != nil {
		logger.Println(err)
	}
}
