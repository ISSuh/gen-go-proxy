# gen-tx-proxy

## Usage
```bash
$ gen-tx-proxy --help
Usage: gen-tx-proxy [--interface-package-name INTERFACE-PACKAGE-NAME] [--interface-package-path INTERFACE-PACKAGE-PATH] --target TARGET [--output OUTPUT] [--package PACKAGE]

Options:
  --interface-package-name INTERFACE-PACKAGE-NAME, -n INTERFACE-PACKAGE-NAME
                         package name of the target interface source code file
  --interface-package-path INTERFACE-PACKAGE-PATH, -l INTERFACE-PACKAGE-PATH
                         package path of the target interface source code file
  --target TARGET, -t TARGET
                         target interface source code file. is required
  --output OUTPUT, -o OUTPUT
                         output file path. default is the same as the target interface source code file
  --package PACKAGE, -p PACKAGE
                         package name of the generated code. default is the same as the target interface source code file
  --help, -h             display this help and exit
```

```bash
# generate proxy to same directory
gen-tx-proxy -t cmd/example/service/foo.go

# generate proxy to other directory and has custom package name
gen-tx-proxy -t ./example/service/bar.go -o ./example/service/proxy -p proxy -n service -l github.com/ISSuh/simple-gen-proxy/example/service
```