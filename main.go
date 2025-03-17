package main

import (
	"fmt"
	"net/http"
)

func main() {

	beans := http.NewServeMux()
	beans.Handle("/",http.FileServer(http.Dir(".")))

	toast := http.Server{
		Handler: beans,
		Addr: ":8080",
	}

	toast.ListenAndServe()

	fmt.Println("hello world")
}