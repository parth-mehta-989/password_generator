# password_generator

Package `password_generator` provides functionality to generate random passwords that meet a specified condition.

To get the package and add it in your project dependencies, use:
```shell
go get -d github.com/parth-mehta-989/password_generator
```

## Usage
Example:
```go
    package main

    import (
        "log"
         gen "github.com/parth-mehta-989/password_generator"
    )

    func main() {
        conditions := gen.Conditions{
            MinUppercase:   1,
            MinLowercase:   1,
            MinNumber:      1,
            MinSpecialChar: 1,
            MinLength:      8,
            MaxLength:      15,
        }
        allowedSpecialChars := "!@#$"
        g := gen.NewGenerator(conditions, &allowedSpecialChars)
        pass, err := g.Generate()
        if err != nil {
            log.Panicf("error generating random password, err:%s", err)
        }

        log.Printf("password: %s, len:%d", *pass, len(*pass))
        /* output:
         xxxx/xx/xx xx:xx:xx password: Ez0@OPWiXo#0, len:12
         Note: your password might differ. :)
        */
    }

```
