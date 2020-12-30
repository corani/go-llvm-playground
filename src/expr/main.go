package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/llvm/llvm-project/llvm/bindings/go/llvm"
)

func parseFile(name string) Ast {
	bs, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalf("error reading source file: %v\n", err)
	}

	p := &Parser{src: string(bs)}
	return p.ParseExpression()
}

func writeLL(llfile *os.File, mod llvm.Module) {
	defer llfile.Close()

	if _, err := llfile.WriteString(mod.String()); err != nil {
		log.Fatalf("error writing LL file: %v\n", err)
	}
}

func emitFile(name string, ast Ast) {
	mod := createModule(ast)

	llfile, err := ioutil.TempFile("", "expr.*.ll")
	if err != nil {
		log.Fatalf("error creating LL file: %v\n", err)
	}
	defer os.Remove(llfile.Name())

	writeLL(llfile, mod)

	// hardcoded call to clang to compile the LL into a binary
	cmd := exec.Command("clang-10", "-o", name, llfile.Name())
	if err := cmd.Run(); err != nil {
		log.Fatalf("error compiling: %v\n", err)
	}
}

func main() {
	var (
		inputName  string
		outputName string
	)

	flag.StringVar(&inputName, "i", "", "input `filename`")
	flag.StringVar(&outputName, "o", "a.out", "output `filename`")
	flag.Parse()

	if inputName == "" {
		flag.PrintDefaults()
		log.Fatalln("no input file specified")
	}

	ast := parseFile(inputName)

	emitFile(outputName, ast)
}
