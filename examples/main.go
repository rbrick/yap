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
		],
		"numbers": [1,2,3,4,5,6,7,8,9,10],
		"elements": [
		  {
		     "test": [2]
		  },
		  {
		    "test": [3] 
		  }
		],
		"largeNumber": 9423233329388648686826386283682368.32
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

	result4, err := yap.Evaluate(`length(where($.books, @.author == "Mary Shelley")) >= 1`, jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Result 4:", result4)

	result5, err := yap.Evaluate(`where($.elements, @.test[0] > 2)`, jsonData)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Result 5:", result5)

	log.Println(
		yap.Evaluate(`largeNumber`, jsonData))
}
