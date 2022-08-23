package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load("balance/infrastructure/config/env/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	secretKey := os.Getenv("SECRET_KEY")

	fmt.Println(s3Bucket, secretKey)

	// now do something with s3 or whatever
}
