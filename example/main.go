package main

import (
	"fmt"
	"github.com/douguohai/gen-id/generator"
)

func main() {
	fmt.Println(generator.GeneratorName())
	fmt.Println(generator.GeneratorIDCart(nil))
	fmt.Println(generator.GeneratorVocationalCertificate())
	fmt.Println(generator.GeneratorPhone())
}
