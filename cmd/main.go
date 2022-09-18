package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Kolyan4ik99/fileServer/internal/app"
)

func main() {
	fileServer := app.NewFileServer("./files_folder", ":8080")

	err := fileServer.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func testFile(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	for i := 0; i < 1000000; i++ {
		file.Write([]byte(fmt.Sprintf("Hello world %d\n", i)))
	}
}
