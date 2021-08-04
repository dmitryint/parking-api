package main

import (
	"log"
	"main/app/api"
	"net/http"
	"os"

	"github.com/apex/gateway"
)

func main() {
	_, isLambdaRuntime := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")

	if isLambdaRuntime {
		log.Fatal(gateway.ListenAndServe(":5000", api.NewRouter()))
	} else {
		log.Fatal(http.ListenAndServe(":5000", api.NewRouter()))
	}
}
