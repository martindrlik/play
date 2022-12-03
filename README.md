# play

Play is a http server that enables creating API by uploading
handler's source code (Go). Compiled as a plugin ready to
handle requests on specified endpoint.

Note that play is my playground for experiments. Without
modifications it should not be taken to production as
there might be security or other issues. Most notably
it compiles and runs uploaded code permitted to do
the same things as play process itself.

## api

### /upload/specified/path

Uploads source code and makes it ready for use on path
`/specified/path`. Note that your source code needs to
have `Main(http.ResponseWriter, *http.Request)` func.

### /analyze/specified/path

If there is an error while deploying source code it should
be investigated on path `/analyze/specified/path`.

### /specified/path

After source code is successfully deployed it serves on
`/specified/path`.

## How it works?

Uploaded content is produced to kafka topic. Error
while uploading is reported directly to client and
also to play's logs (stdout) and metrics. Play
is suscribed to consume that topic's messages.
Once message is pulled it builds source code
provided by message as a Go plugin and lookups
Main func. Error during this process is available
on path `/analyze/specified/path`. It should be
always consulted.

If no error `/specified/path` can be requested.

Hurray! 🎉

### Replicas

In order to improve availability and to handle
more requests multiple play servers can be
deployed. Client is then usually requesting
load balancer which forwards client request
to by some algorithm choosen replica. Round
robin or more sophisticated algorithm can
be used. All replicas are the same so it
does not matter which one is picked for
the job.

#### Replicas can be added or removed

If higher load (more clients) is expected or otherwise
play servers can be added or removed as needed.

If a new replica is added then it starts processing
all kafka messages from specified topic, from
the oldest one to the most recent one. It
means that new replica will eventually
catch up with others.

Note that play has bug here! New replica will start
serving even though it did not catch up with
others yet.

### Result consistency

Note that play has eventually consistent model.
Client only waits for message to be send (produced).
Client does not wait for message to be processed,
source code compiled and so on. This happens
asynchronously. It means that result of requesting
`/specified/path` or `/analyze/specified/path`
right after uploading may not be available.

With multiple replicas and also depends on deployed
load balancing strategy the result might be
inconsistent, but eventually become consistent
as all replicas complete processing message
and compilation.

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

Star kafka broker and play server. Upload `hello.go` file.

```zsh
% curl --data-binary @hello.go http://localhost:8085/upload/say/hello
```

Note that `/upload/{specified/path}` is `/upload/say/hello`. `/say/hello` is going to be used as a path for calling that "uploaded api".

```zsh
% curl http://localhost:8085/say/hello\?name=Gopher
Hello, Gopher!
```
