# gen-tx-proxy

## Description
Simple Gen Proxy is a tool for generating proxy source code files for interfaces within Go source code files. It simplifies the creation and management of proxy code, making it easier for developers to implement proxy patterns in their Go applications.

## Installation
To install Simple Gen Proxy, use the following command:

```bash
go get github.com/ISSuh/gen-tx-proxy
```

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

## example

implement interface of service layer

```go
package service

type Foo interface {
    // add the @transactionl comments to the method you want to apply the transaction
    // @transactionl
    Create(email string) (int64, error)

    // If there is no comment @transactionl,
    // the method calls the function immediately without the transaction being applied
    Find(id int64) (string, error)
}
```



## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.