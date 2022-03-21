package main

import (
	"nubes/sum/db"
	"nubes/sum/router"
	"nubes/sum/utils"
)

func main() {
	utils.PrepareKeys()
	db.InitDatabase()

	router := router.New()
	router.Run("localhost:8080")
}
