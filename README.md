# Go LLVM Playground

Here I'm collecting some sample code I wrote while playing around with the Golang bindings for
LLVM, as they may help others get started.

## Setup

I'm using Ubuntu 18.04 and Go 1.14.9 on a system that I've used for development for some time
already, so I've no idea which dependencies are required, as most of them were installed already.

1. Get the LLVM bindings:
   `go get -d github.com/llvm/llvm-project/llvm/bindings/go/llvm`

   Note: This pulls down the entire LLVM project, which is rather large. This will take a while!

2. Build the LLVM bindings:
   ```bash
   cd $GOPATH/src/github.com/llvm/llvm-project/llvm/bindings/go
   ./build.sh -DLLVM_TARGETS_TO_BUILD=host
   ```

   Note: This builds the entire LLVM project, so this will take a while (and use a few GB of disk)

The above will generate a `llvm_config.go` file in the `llvm` package with the CGO flags baked in.

Note: I had to install cmake from the [Kitware APT Repository](https://apt.kitware.com), as the
version in the Ubuntu repository was too old.

Note: I had to *uninstall* ninja-build, as I couldn't get this to build the project without errors
about missing libraries.

## Examples

- [src/add/main.go](src/add/main.go)

  Export a simple Go function to the C namespace, generate a new module using LLVM that defines
  a `main` function with two local variables. Invoke the Go function from the `main` function and
  return the result.

  Use the LLVM JIT compiler to execute the `main` function from within Go and print the result.

- [src/hello/main.go](src/hello/main.go)

  Generate a new module using LLVM that defines a global constant for "Hello, world!" and a
  `main` function that prints the string using `puts`. Write the LL, ASM and OBJ code for this
  module to the current working directory, then link the object file into an executable.

  I have not found a way to generate an executable directly from the LLVM bindings, so I had to
  call out to `gcc` to do this.

## Resources

- The "Add" example borrows heavily from [An introduction to LLVM in Go](https://felixangell.com/blogs/an-introduction-to-llvm-in-go).
- The "Hello" example takes some inspiration from [llvm-hello-world](https://github.com/dfellis/llvm-hello-world).
