package main

import (
	"io/ioutil"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func createModule() *ir.Module {
	i32_t := types.I32

	mod := ir.NewModule()

	// static char* hello = "Hello, world!\n"
	hello := mod.NewGlobalDef("hello", constant.NewCharArrayFromString("Hello, world!\n"))

	// extern int puts(...)
	puts := mod.NewFunc("puts", i32_t)
	puts.Sig.Variadic = true

	// int main(void)
	main := mod.NewFunc("main", i32_t)
	block := main.NewBlock("entry")

	// puts(hello)
	block.NewCall(puts, hello)

	// return 0
	block.NewRet(constant.NewInt(i32_t, 0))

	return mod
}

func writeLL(mod *ir.Module) {
	if err := ioutil.WriteFile(
		"hello.ll", []byte(mod.String()), 0644,
	); err != nil {
		panic(err)
	}
}

func main() {
	mod := createModule()

	writeLL(mod)
}
