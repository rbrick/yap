package main

import (
	"log"

	"github.com/rbrick/yap"
)

func main() {
	jsonData := `{
		"books": [
			{"name": "Frankenstein", "author": "Mary Shelley"},
			{"name": "1984", "author": "George Orwell"},
			{"name": "Project Hail Mary 🌌", "author": "Andy Weir"}
		]
	}`

	result1, err := yap.Evaluate(`$.books[0].name`, jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Result 1:", result1)

	result2, err := yap.Evaluate(`length($.books) >= 2`, jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Result 2:", result2)

	// works with unicode too!
	result3, err := yap.Evaluate(`equals($.books[2].name, "Project Hail Mary 🌌")`, jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Result 3:", result3)
}
