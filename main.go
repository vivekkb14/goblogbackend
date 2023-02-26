package main

import (
	"fmt"

	"github.com/vivekkb14/goblogbackend/dbops"
	"github.com/vivekkb14/goblogbackend/httpserver"
)

func main() {
	dbops.GlobalDatabase = new(dbops.DataBase)
	err := dbops.GlobalDatabase.InitialiseDatabaeServer()
	if err != nil {
		fmt.Println("Error in initialising database")
		return
	}
	fmt.Println("Connected!")
	dbops.GlobalDatabase.CreateProductTable()
	httpserver.CreateHttpServer()
}
