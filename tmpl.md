# the Tmpl API

## Terminate tmpl

```go
{{- if false }}
{{- $t := .Terminate -}}
{{- end -}}
```

## Glob

### ToLowerCamelCase

convert string to lower camel case

-   string

`{{- ToLowerCamelCase "HelloWorld" }}`

> helloWorld

### ToUpperCamelCase

convert string to upper camel case

-   string

`{{- ToUpperCamelCase "helloWorld" }}`

> HelloWorld

### SnakeStringToBigHump

convert snake string to big hump

-   string

`{{- ToUpperCamelCase "hello_world" }}`

> HelloWorld

### Contains

check if string contains substring or a string in a string array

-   bool

`{{- Contains "hello_world" "hello" }}`

> true

`{{- Contains []string{"hello"} "hello_world" }}`

> false

## Entity

### ListFieldByTag

get all field by tag

-   []Field

`{{- range (.ListFieldByTag "gomate:getter") }}`
