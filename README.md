# play

Play is my pet project where I play with web server architecture. It is a http server that allows to upload go file which is then built as plugin and runs by calling url path returned by `/upload` api. Uploaded go file is expected to have `func Main(http.ResponseWriter, *http.Request)` in main package. The function will be called by `/run/{token}` api.

## Example

Start web server and upload following go file by using `/upload` api.

```go
package main

import (
    "fmt"
    "net/http"
)

func Main(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "Hello, world!")
}
```

Response is a path to be used for calling `Main` function of that uploaded file. Response example:

```text
/run/rt
```

Finally call `/run/{token}` api to execute that `Main` function:

```zsh
% curl http://:8080/run/rt
Hello, world!
```
