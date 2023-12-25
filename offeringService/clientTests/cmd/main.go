package main

import (
	"Go_project/offeringService/clientTests"
	"Go_project/offeringService/clientTests/api"
)

func main() {
	url := "http://localhost:8080"
	currentClient := client.NewClient(url)
	api.TestRequests(currentClient)
}
