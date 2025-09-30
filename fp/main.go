package main

import "fmt"

// dev-1
func hello() {
	fmt.Println("Hello")
}

// dev-2
func hi() {
	fmt.Println("Hi")
}

func hey() {
	fmt.Println("Hey")
}

// design issues
// - code ( emjo ) duplication / scattering
// - code tangling / tight coupling. All functions are tightly coupled with emoji printing logic

// solution: High Order Functions (HOFs)

// HOF: A function that either takes a function as argument or returns a function as result
func emojiDecorator(f func(), emoji string) func() {
	return func() {
		fmt.Print(emoji, " ")
		f()
		fmt.Print(" ", emoji)
	}
}

func main() {

	hello()
	emojiHello := emojiDecorator(hello, "👋")
	emojiHello()

	hi()
	emojiHi := emojiDecorator(hi, "🖐️")
	emojiHi()

	hey()
	emojiHey := emojiDecorator(hey, "🤚")
	emojiHey()

}
