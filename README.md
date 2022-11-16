# play

Play is a http server that enables creating API by uploading go file.

## api

### /upload/{specified/path}

Uploads Go file and makes it ready for use on path /specified/path.

If the {specified/path} starts with `notify/` calling this API will also produced path to configured kafka topic.

If the {specified/path} starts with `subscribe/` API will be called by any message pulled from configured kafka topic (`-kafka-consumer-topic`).

## Example

Create a go file called `hello.go` with following content.

```go
package main

import (
	"fmt"
	"net/http"
)

func Main(rw http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	fmt.Fprintf(rw, "Hello, %s!\n", name)
}
```

Star play server. Upload `hello.go` file.

```zsh
% curl --data-binary @hello.go http://localhost:8085/upload/say/hello
```

Note that `/upload/{specified/path}` is `/upload/say/hello`. `/say/hello` is going to be used as a path for calling that "uploaded api".

```zsh
% curl http://localhost:8085/say/hello\?name=Gopher
Hello, Gopher!
```
