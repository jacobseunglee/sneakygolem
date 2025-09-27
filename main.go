package main

import (
	"fmt"
	files "sneakygolem/utils"
)

func main() {
	fmt.Println("hello")
	content := files.ReadContent("test.pdf")
	fmt.Println(content)
	files.WriteContent(content, "hello")
}
