package main

import "fmt"

func main() {
	//экземпляры структурок
	plain := PlainText{text: "happy"}
	bolt := BoltText{text: "lama"}
	code := CodeText{text: "no"}
	italic := ItalicText{text: "drama"}

	// вывод формата
	fmt.Printf("Plain: %s\n", plain.Format())
	fmt.Printf("Bolt: %s\n", bolt.Format())
	fmt.Printf("Code: %s\n", code.Format())
	fmt.Printf("Italic: %s\n", italic.Format())

	// Цепочка вариант 1
	chF1 := ChainFormater_1{text: "happy lama"}
	chF1.AddFormatter(plain)
	chF1.AddFormatter(bolt)
	chF1.AddFormatter(code)
	chF1.AddFormatter(italic)
	chF1.AddFormatter(italic)
	chF1.AddFormatter(code)
	chF1.AddFormatter(bolt)

	// Цепочка 2
	chF2 := ChainFormater_2{}
	chF2.AddFormatter(plain)
	chF2.AddFormatter(bolt)
	chF2.AddFormatter(code)
	chF2.AddFormatter(italic)
	// Вывод
	fmt.Printf("ChainFormater_1: %s\n", chF1.Format())
	fmt.Printf("ChainFormater_2: %s\n", chF2.Format())

}
