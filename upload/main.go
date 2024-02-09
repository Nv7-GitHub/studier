package main

import (
	"encoding/json"
	"os"
)

const file = "../sets/apush/ch/ch21.st"

func main() {
	qs := ParseFile(file)

	// Create
	res, _ := json.Marshal(qs)
	f, _ := os.Create("data.json")
	f.Write(res)
}
