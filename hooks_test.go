package main

import (
	"log"
	"os"
)

func setup() {
	cfg := initApiConfig()
	err := deleteAllUsersAndPosts(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func tearDown() {
	cfg := initApiConfig()
	err := deleteAllUsersAndPosts(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
