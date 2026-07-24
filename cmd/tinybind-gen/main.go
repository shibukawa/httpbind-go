package main

import "github.com/shibukawa/tinybind-go/generator"

func main() {
	generator.Main(generator.MustCommandSet(generator.GenerateCommand(generator.DefaultOptions())))
}
