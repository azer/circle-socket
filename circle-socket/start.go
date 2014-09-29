package main

import (
	"fmt"
	"github.com/azer/circle-socket"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	_ = godotenv.Load()
	circle.CreateFlickrClient()
	circle.CreateDBConn("./data-flickr")
	circle.Start(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
