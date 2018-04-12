package main

import (
	"whatapp/fs"
	"log"
)

const (
	dataFolder = "D:\\Users\\81092610\\Desktop\\data"
)

func main() {
	var err error
	err = fs.Descompacta(dataFolder)
	if err != nil {
		log.Fatal(err)
	}
	err = fs.MesclaTxt(dataFolder)
	if err != nil {
		log.Fatal(err)
	}
}
