# go-conf

A straightforward, distributed, file-based configuration solution for Go.

## Install

```bash
go get -u github.com/acrazing/go-conf
```

## Usage

1. In any module, register config node and read it with `conf.Registry`:

    ```go
    package my

    import conf "github.com/acrazing/go-conf"

    type Config struct {
        User string `json:"user"`
    }

    func init() {
        conf.Register[Config]("my.path")
    }

    func New(r *conf.Registry) {
        config, err := conf.Get[Config](r)
        if err != nil {
            panic(err)
        }
        println(config.User)
    }
    ```

2. Load configuration file and get conf registry in your app entry:

    ```go
    package main

    import conf "github.com/acrazing/go-conf"

    func main() {
        r, err := conf.Load("config.json")
        if err != nil {
            panic(err.Error())
        }
        my.New(r)
    }
    ```

3. Helper: print default config file structure:

    ```go
    package main

    import conf "github.com/acrazing/go-conf"
    import "encoding/json"

    func main() {
        c := conf.Default()
        jv, _ := json.MarshalIndent(c, "", "  ")
        println(string(jv))
    }
    ```
