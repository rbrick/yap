package yap

type Parser struct {
}

func Parse(input string) {
	// two libraries we want to conjoin
	// jsonpath & conditions

	// "json(key.field).field2" -> resolver, no binop -> parses key.field as JSON by calling JSON function, then returns field2
	// "json(key.field) == value" -> resolver, binop -> parses key.field as JSON by calling JSON function, then compares to value
	// binops always return booleans
	// functions can return any type
	// functions should be able to be defined
}
