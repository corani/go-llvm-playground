package main

// extern int AddInts(int, int);
import "C"

import (
	"fmt"

	"github.com/llvm/llvm-project/llvm/bindings/go/llvm"
)

//export AddInts
func AddInts(a, b C.int) C.int {
	return a + b
}

func createModule() llvm.Module {
	builder := llvm.NewBuilder()

	mod := llvm.NewModule("my_module")
	mod.SetTarget(llvm.DefaultTargetTriple())

	// types
	i32_t := llvm.Int32Type()

	// extern int AddInts(int, int)
	addIntsType := llvm.FunctionType(i32_t, []llvm.Type{i32_t, i32_t}, false)
	addInts := llvm.AddFunction(mod, "AddInts", addIntsType)

	// int main(void)
	main := llvm.FunctionType(i32_t, []llvm.Type{}, false)
	llvm.AddFunction(mod, "main", main)
	block := llvm.AddBasicBlock(mod.NamedFunction("main"), "entry")
	builder.SetInsertPoint(block, block.FirstInstruction())

	// int a = 32
	a := builder.CreateAlloca(i32_t, "a")
	builder.CreateStore(llvm.ConstInt(i32_t, 32, false), a)

	// int b = 16
	b := builder.CreateAlloca(i32_t, "b")
	builder.CreateStore(llvm.ConstInt(i32_t, 16, false), b)

	// res = AddInts(a, b)
	aval := builder.CreateLoad(a, "a_val")
	bval := builder.CreateLoad(b, "b_val")
	call := builder.CreateCall(addInts, []llvm.Value{aval, bval}, "res")

	// return res
	builder.CreateRet(call)

	return mod
}

func verifyModule(mod llvm.Module) {
	if err := llvm.VerifyModule(
		mod, llvm.VerifierFailureAction(llvm.ReturnStatusAction),
	); err != nil {
		panic(err)
	}
}

func createEngine(mod llvm.Module) llvm.ExecutionEngine {
	engine, err := llvm.NewExecutionEngine(mod)
	if err != nil {
		panic(err)
	}

	return engine
}

func run(engine llvm.ExecutionEngine, mod llvm.Module) uint64 {
	res := engine.RunFunction(mod.NamedFunction("main"), []llvm.GenericValue{})

	return res.Int(false)
}

func main() {
	llvm.LinkInMCJIT()
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()

	mod := createModule()
	verifyModule(mod)

	engine := createEngine(mod)

	// dump the generated LL to stdout
	mod.Dump()

	res := run(engine, mod)
	fmt.Printf("result: %d\n", res)
}
