# the Tmpl API

## Glob

### ToLowerCamelCase

convert string to lower camel case

return: string

{{- ToLowerCamelCase "HelloWorld" }}

output: helloWorld

### ToUpperCamelCase

convert string to upper camel case

return: string

{{- ToUpperCamelCase "helloWorld" }}

output: HelloWorld

### SnakeStringToBigHump

convert snake string to big hump

return: string

{{- ToUpperCamelCase "hello_world" }}

output: HelloWorld

## Entity

### ListFieldByTag

get all field by tag

return: []Field

{{- range (.ListFieldByTag "gomate:getter") }}
