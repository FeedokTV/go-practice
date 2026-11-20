package methodsets

// Covers the information from below:
// https://go.dev/ref/spec#Struct_types
// https://go.dev/ref/spec#Method_sets
// https://go.dev/ref/spec#Calls
// https://go.dev/ref/spec#Address_operators
// https://go.dev/doc/effective_go#embedding

type Speaker interface {
	Speak() string
}

// ----- Dog -----

type Dog struct{}

func (d Dog) Speak() string {
	return "dog!"
}

// ----- Pdog -----

type PDog struct{}

func (pd *PDog) Speak() string {
	return "pdog"
}

// ----- Embeds -----

type EmbedDog struct {
	Dog
}

type EmbedPDog struct {
	PDog
}

type EmbedPointPDog struct {
	*PDog
}

type EmbedPointDog struct {
	*Dog
}

// ----- Actual Tests -----

var _ Speaker = Dog{} // Compiles: regular use (T expose all func(t T))

var _ Speaker = &Dog{} // Compiles: *T expose all func(pt *T) and func(t T) as well

// var _ Speaker = PDog{}  Not Compiles: here we have func(pt *T) for PDog and PDog won't expose it, but *PDog - yes (see example above)
// Compiler Error: PDog does not implement Speaker (method Speak has pointer receiver)

var _ Speaker = &PDog{} // Compiles: Fix of version above, *T expose to both types

var _ Speaker = EmbedDog{} // Compiles: regular use (T embedded in S - regular func(t T) will be able in S)

var _ Speaker = &EmbedDog{} // Compiles: T embedded in *S - func(t T) will be able in S and func(pt *T) as well

// var _ Speaker = EmbedPDog{}  Not Compiles:  T embedded in S not promoting the func(pt *T) funcs
// Compiler Error: EmbedPDog does not implement Speaker (method Speak has pointer receiver)

var _ Speaker = &EmbedPDog{} // Compiles: T embed in *S expose functions with receiver T and *T and S and *S

var _ Speaker = EmbedPointDog{} // Compiles: *T embed in S, funcs for T, *T and S available

var _ Speaker = &EmbedPointDog{} // Compiles: *T embed in *S, funcs for T, *T and S, *S available

var _ Speaker = EmbedPointPDog{} // Compiles: *T embed in S, funcs for T, *T and S available

var _ Speaker = &EmbedPointPDog{} // Compiles: *T embed in *S, funcs for T, *T and S, *S available
