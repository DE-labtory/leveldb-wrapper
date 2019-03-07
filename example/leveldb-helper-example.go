package main

import (
	"fmt"
	"log"
	"os"

	leveldbwrapper "github.com/DE-labtory/leveldb-wrapper"
)

func main() {

	path := "./leveldb"
	defer os.RemoveAll(path)
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()

	err := db.Put([]byte("20164403"), []byte("JUN"), true)

	if err != nil {
		log.Fatalf("error occured [%s]", err.Error())
	}

	name, _ := db.Get([]byte("20164403"))

	fmt.Printf("%s", name)
}
