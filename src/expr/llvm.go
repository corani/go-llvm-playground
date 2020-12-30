package main

import "github.com/llvm/llvm-project/llvm/bindings/go/llvm"

func createModule(ast Ast) llvm.Module {
	builder := llvm.NewBuilder()

	mod := llvm.NewModule("main")
	mod.SetTarget(llvm.DefaultTargetTriple())

	// types
	i32_t := llvm.Int32Type()

	// extern printf(...)
	printfType := llvm.FunctionType(llvm.VoidType(), []llvm.Type{}, true)
	printf := llvm.AddFunction(mod, "printf", printfType)

	// int main(void)
	main := llvm.FunctionType(i32_t, []llvm.Type{}, false)
	llvm.AddFunction(mod, "main", main)
	block := llvm.AddBasicBlock(mod.NamedFunction("main"), "entry")
	builder.SetInsertPoint(block, block.FirstInstruction())

	// int lhs = xxx
	lhs := builder.CreateAlloca(i32_t, "")
	builder.CreateStore(llvm.ConstInt(i32_t, ast.lhs, false), lhs)

	// int rhs = yyy
	rhs := builder.CreateAlloca(i32_t, "")
	builder.CreateStore(llvm.ConstInt(i32_t, ast.rhs, false), rhs)

	lhs_val := builder.CreateLoad(lhs, "")
	rhs_val := builder.CreateLoad(rhs, "")

	var res llvm.Value

	switch ast.oper {
	case "+":
		// res = lhs + rhs
		res = builder.CreateAdd(lhs_val, rhs_val, "")
	case "-":
		// res = lhs - rhs
		res = builder.CreateSub(lhs_val, rhs_val, "")
	}

	// printf("result: %d\n", res)
	pattern := builder.CreateGlobalStringPtr("result: %d\n", "pattern")
	builder.CreateCall(printf, []llvm.Value{pattern, res}, "")

	// return res
	builder.CreateRet(res)

	return mod
}
