package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		switch os.Args[1] {
		case "read":
			{
				if bytes, err := os.ReadFile(os.Args[2]); err != nil {
					fmt.Printf("could not be read %s: %s\n", os.Args[2], err)
				} else {
					fmt.Println(string(bytes))
				}
			}
		case "create":
			{
				if file, err := os.Create(os.Args[2]); err != nil {
					fmt.Printf("could not be create %s:%s\n", os.Args[2], err)
				} else {
					defer file.Close()
					fmt.Printf("File %s created", os.Args[2])
				}
			}
		case "delete":
			{
				if err := os.Remove(os.Args[2]); err != nil {
					fmt.Printf("could not be delete %s:%s\n", os.Args[2], err)
				} else {
					fmt.Printf("File %s deleted", os.Args[2])
				}
			}
		}
	}

}
