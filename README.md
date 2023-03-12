```go
import "go.szostok.io/bye"
```

Go package providing shutdown manager. 

## Features

- Register shutdown via empty functions, error functions or service with `Shutdown() error` method
- Optionally set shutdown timeout
- All registered functions/services are shutdown in parallel 
- All shutdown errors are aggregated and returned as a single one

## Installation

```shell
go get go.szostok.io/bye
```

## Usage

```go
shutdown := bye.NewParentService(bye.WithTimeout(30 * time.Second))

shutdown.Register(bye.Func(func() {
    fmt.Println("Closing non error function call")
}))

shutdown.Register(bye.ErrFunc(func() error {
    fmt.Println("Closing error function call")
    return errors.New("I don't want to quit!")
}))

shutdown.Register(&exampleService{})

<- close

fmt.Println("Shutting down the application...")

err := shutdown.Shutdown()
fmt.Println(err)
```
