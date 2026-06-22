package main
import("fmt")

func bro() {
	fmt.Println("Hi bro!")
}
func main() {
	fmt.Println("Hi guys")
	bro()
	what()
}
func what() {
	fmt.Println("What?")
}