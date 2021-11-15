# play

Play is a http server that allows you to create api by uploading go file. Go file is then built as a plugin and it is ready to be called on specified path.

## api

### /upload/{specified/path}

Uploads go file and makes it ready for calling. After /upload there should be a path on which you can then call that "uploaded api".

## Example

Start play web server. Create a go file called `hello.go` with following or similar content. Note that it is only expected to be main package with Main function.

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

Upload that file. We can use `curl` for instance.

```zsh
% curl --data-binary @hello.go http://localhost:8085/upload/say/hello
```

Note that `/upload/{specified/path}` is `/upload/say/hello`. `/say/hello` is going to be path for calling that "uploaded api".

```zsh
% curl "http://localhost:8085/say/hello?name=Gopher"
Hello, Gopher!
```
