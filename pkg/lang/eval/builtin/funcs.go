package builtin

type GenericFunc func(args []interface{}) (interface{}, error)

var Functions = map[string]GenericFunc{
	// math
	"+": sumFunction,
	"-": minusFunction,
	"*": multiplyFunction,
	"/": divideFunction,

	// strings
	"concat": concatFunction,
	"split":  splitFunction,

	// lists
	"len": lenFunction,

	// logic
	"and": andFunction,
	"or":  orFunction,
}
