package main

import (
	"fmt"
	"os"

	leveldbwrapper "github.com/DE-labtory/leveldb-wrapper"
)

func main() {

	path := "./leveldb"
	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
	defer os.RemoveAll(path)

	studentDB := dbProvider.GetDBHandle("Student")
	studentDB.Put([]byte("20164403"), []byte("JUN"), true)

	name, _ := studentDB.Get([]byte("20164403"))

	fmt.Printf("%s", name)
}
