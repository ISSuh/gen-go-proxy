# proxy example

This example guides you on how to generate and use proxy code with simple-gen-proxy.

## Usage

### implement interface and middleware

```go
package service

// implement target interface
type Foo interface {
  // use @proxy comment if you want to generate proxy code
  // @proxy
  Logic(needEmitErr bool) (string, error)
}
```

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

### generate proxy code

```bash
$ simple-gen-proxy -t service
```

### use generated proxy code

```go
// create instence and regist middleware
func main() {
  target := service.NewFoo()
  proxy := proxy.NewFooProxy(target, Wrapped, Before, After)

  if val, err := proxy.Logic(false); err != nil {
    fmt.Println("err: ", err)
  } else {
    fmt.Println("val: ", val)
  }

  fmt.Println()

  if val, err := proxy.Logic(true); err != nil {
    fmt.Println("err: ", err)
  } else {
    fmt.Println("val: ", val)
  }
}
```

```bash
$ go run main.go
[Wrapped] before
[Before] before
[Foo] logic
[After] after
[Wrapped] after
val:  foo logic

[Wrapped] before
[Before] before
[Foo] logic
[After] err occurred. err : emit error
[After] after
[Wrapped] err occurred. err : emit error
[Wrapped] after
err:  emit error
```
