package main

import (
	"io/ioutil"
	"os/exec"

	"github.com/llvm/llvm-project/llvm/bindings/go/llvm"
)

func createModule() llvm.Module {
	builder := llvm.NewBuilder()

	mod := llvm.NewModule("my_module")
	mod.SetTarget(llvm.DefaultTargetTriple())

	// types
	i32_t := llvm.Int32Type()
	charptr_t := llvm.PointerType(llvm.Int8Type(), 0)

	// extern int puts(char*)
	putsType := llvm.FunctionType(i32_t, []llvm.Type{charptr_t}, false)
	puts := llvm.AddFunction(mod, "puts", putsType)

	// int main(void)
	main := llvm.FunctionType(i32_t, []llvm.Type{}, false)
	llvm.AddFunction(mod, "main", main)
	block := llvm.AddBasicBlock(mod.NamedFunction("main"), "entry")
	builder.SetInsertPoint(block, block.FirstInstruction())

	// static char* hello = "Hello, world!\n"
	hello := builder.CreateGlobalStringPtr("Hello, world!\n", "hello")

	// puts(hello)
	builder.CreateCall(puts, []llvm.Value{hello}, "")

	// return 0
	builder.CreateRet(llvm.ConstInt(i32_t, 0, false))

	return mod
}

func verifyModule(mod llvm.Module) {
	if err := llvm.VerifyModule(
		mod, llvm.VerifierFailureAction(llvm.ReturnStatusAction),
	); err != nil {
		panic(err)
	}
}

func createMachine(mod llvm.Module) llvm.TargetMachine {
	target, err := llvm.GetTargetFromTriple(mod.Target())
	if err != nil {
		panic(err)
	}

	return target.CreateTargetMachine(mod.Target(), "", "",
		llvm.CodeGenLevelDefault,
		llvm.RelocPIC,
		llvm.CodeModelDefault)
}

func writeLL(mod llvm.Module) {
	if err := ioutil.WriteFile(
		"hello.ll", []byte(mod.String()), 0644,
	); err != nil {
		panic(err)
	}
}

func writeASM(machine llvm.TargetMachine, mod llvm.Module) {
	buffer, err := machine.EmitToMemoryBuffer(mod, llvm.AssemblyFile)
	if err != nil {
		panic(err)
	}
	defer buffer.Dispose()

	if err := ioutil.WriteFile(
		"hello.s", buffer.Bytes(), 0644,
	); err != nil {
		panic(err)
	}
}

func writeOBJ(machine llvm.TargetMachine, mod llvm.Module) {
	buffer, err := machine.EmitToMemoryBuffer(mod, llvm.ObjectFile)
	if err != nil {
		panic(err)
	}
	defer buffer.Dispose()

	if err := ioutil.WriteFile(
		"hello.o", buffer.Bytes(), 0644,
	); err != nil {
		panic(err)
	}
}

func compile() {
	cmd := exec.Command("gcc", "-o", "hello", "hello.o")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func main() {
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()

	mod := createModule()
	verifyModule(mod)

	machine := createMachine(mod)

	writeLL(mod)
	writeASM(machine, mod)
	writeOBJ(machine, mod)
	compile()
}
