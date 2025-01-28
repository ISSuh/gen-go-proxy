# gen-go-proxy

## Description

gen-go-proxy is a tool for generating proxy source code files for interfaces within Go source code files. It simplifies the creation and management of proxy code, making it easier for developers to implement proxy patterns in their Go applications.

Implement basic proxy patterns through generated proxy codes.

You can also display a dedicated annotation(comment) in the method of the interface to which you want to apply the proxy, allowing the user to register the middleware, specifying which functions to call in advance during proxy calls.

## Installation

| not for now

To install Simple Gen Proxy, use the following command:

```bash
go get github.com/ISSuh/gen-go-proxy
```

## Usage

```bash
$ gen-go-proxy --help
Usage: gen-go-proxy [--interface-package-name INTERFACE-PACKAGE-NAME] [--interface-package-path INTERFACE-PACKAGE-PATH] --target TARGET [--output OUTPUT] [--package PACKAGE] [--use-tx-middleware]

Options:
  --interface-package-name INTERFACE-PACKAGE-NAME, -n INTERFACE-PACKAGE-NAME
                         package name of the target interface source code file
  --interface-package-path INTERFACE-PACKAGE-PATH, -l INTERFACE-PACKAGE-PATH
                         package path of the target interface source code file
  --target TARGET, -t TARGET
                         target directory path of the interface source code file. is required
  --output OUTPUT, -o OUTPUT
                         output file path. default is the same as the target interface source code file
  --package PACKAGE, -p PACKAGE
                         package name of the generated code. default is the same as the target interface source code file
  --use-tx-middleware, -x
                         generate transaction middleware. default is false
  --help, -h             display this help and exit
```

```bash
# generate proxy to same directory
$ gen-go-proxy -t cmd/example/service

# generate proxy to other directory and has custom package name
$ gen-go-proxy -t ./example/service \
              -o ./example/service/proxy \
              -p proxy \
              -n service \
              -l github.com/ISSuh/gen-go-proxy/example/service

# generate proxy and pre implement transaction middleware
# to other directory and has custom package name
$ gen-go-proxy -t ./example/service \
              -o ./example/service/proxy \
              -p proxy \
              -n service \
              -l github.com/ISSuh/gen-go-proxy/example/service \
              -x
```

## Example

implement interface and adjust user custom annotation for method

```go
package service

type Foo interface {
  // use @{annotation name} comment if you want to generate proxy code
  // @proxy
  A() (string, error)

  // also support multiple annotation
  // proxy middleware runs in order of annotation
  // @custom1
  // @custom2
  B(c context.Context) int
}
```

generate proxy code

```bash
# will create {interface source file name}_proxy.go on same dir
$ gen-go-proxy -t path/target/dir
```


implement user defined middleware for proxy

```go
// implement middleware
func Wrapped(next func(c context.Context) error) func(context.Context) error {
  return func(c context.Context) error {
    fmt.Println("[Wrapped] before")

    // run next middleware or target logic
    err := next(c)
    if err != nil {
      fmt.Printf("[Wrapped] err occurred. err : %s\n", err)
    }

    fmt.Println("[Wrapped] after")
    return err
  }
}

func Before(next func(c context.Context) error) func(context.Context) error {
  return func(c context.Context) error {
    fmt.Println("[Before] before")

    // run next middleware or target logic
    return next(c)
  }
}

func After(next func(c context.Context) error) func(context.Context) error {
  return func(c context.Context) error {
    // run next middleware or target logic
    err := next(c)
    if err != nil {
      fmt.Printf("[After] err occurred. err : %s\n", err)
    }

    fmt.Println("[After] after")
    return err
  }
}
```

create target instence and adjust middleware to proxy instence.

Called in order of the registered middleware when registering multiple middleware in the annotation.

Also if multiple annotations are applied to one method, middleware is called in the order of declared annotations.

```go
  target := service.NewFoo()

  // middleware by annotation
  // key: annotation name
  // value: middleware list
  // can use middleware helper type
  // or raw type map[string][]func(func(context.Context) error) func(context.Context) error
  //
  //  m := map[string][]func(func(context.Context) error) func(context.Context) error{
  //    "proxy":   {Wrapped, Before, After},
  //    "custom1": {Wrapped},
  //    "custom2": {Before, After},
  //  }
  m := service.FooProxyMiddlewareByAnnotation{
    "proxy":   {Wrapped, Before, After},
    "custom1": {Wrapped},
    "custom2": {Before, After},
  }

  // if use middleware helper type,
  // should call helper.To() when create proxy
  proxy := service.NewFooProxy(target, m.To())

  proxy.A()
  proxy.B()
```

### Transaction middleware

This project began with the idea of a better way to handle DB transactions in Golang.

In the commonly used layered architecture, there was a problem with importing a service layer but direct DB-related package using a wrapper function for DB transactions to a service layer with business logic.

I was inspired by Java Spring's transaction processing while thinking about these issues, so I thought about transaction processing through proxies.

Through a separate proxy, the service layer is implemented to form a dependency for business logic only, and the actual DB transaction processing is handled by applying the generated proxy code and middleware.

When a transaction invokes a several sub-service code to which it is applied, it is implemented to first create a transaction and store it in context so that it is processed as a single transaction.

#### Usage

implement interface and adjust user custom annotation for method

The method using transaction mildleware must have **context.Context** in its parameters by default and **error** as a return value.

```go
package service

type Far interface {
  // use @{annotation name} comment if you want to generate proxy code
  // @transactional
  Logic(c context.Context, item dto.Item) (dto.Item, error)
}

type Bar interface {
  // use @{annotation name} comment if you want to generate proxy code
  // @transactional
  Logic(c context.Context, item dto.Item) (dto.Item, error)
}
```

generate proxy code

```bash
# will create
# {interface source file name}_proxy.go and proxy_middleware_tx.go on path/target/dir/proxy
$ gen-go-proxy -t path/target/dir \
              -o path/target/dir/proxy \
              -p proxy \
              -n service \
              -l github.com/ISSuh/path/to/service \
              -x
```

regist tx middleware to annotation

```go
  farTarget := service.NewFoo()  
  barTarget := service.NewBar()

  // This transaction factory function is a function that
  // returns the transaction instance of the DB package that you use.
  txFatory := func() (proxy.Transaction, error) {
    tx := infraorm.NewGORMTransaction(db)
    return tx, nil
  }

  // create transaction middleware and regist
  txMiddleware := proxy.TxMiddleware(txFatory)
  m := map[string][]func(func(context.Context) error) func(context.Context) error{
    "transactional": {txMiddleware},
  }

  foo := service.NewFooProxy(farTarget, m)
  bar := service.NewBarProxy(barTarget, m)
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.