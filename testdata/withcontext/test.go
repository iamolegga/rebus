package withcontext

// Command1 .
// +rebus:out=./bus
type Command1 struct {
}

// Query1 .
// +rebus:out=./bus
type Query1 struct {
	Foo string
}

type Query1Result struct {
	Bar int
}
