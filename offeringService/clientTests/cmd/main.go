package main

import (
	client "Go_project/clientTests"
	"Go_project/clientTests/api"
)

func main() {
	url := "http://localhost:8080"
	currentClient := client.NewClient(url)
	api.TestRequests(currentClient)
}
