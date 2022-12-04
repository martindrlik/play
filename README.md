# play

Play is a HTTP server that allows you to create API
endpoints by uploading a GO file. It is built
as a plugin which is then ready to handle requests
on the specified path.

Note that "play" is my playground for experiments. Without
modifications it should not be taken to production
(as there might be security or other issues). Most notably
it compiles and runs uploaded code permitted to do
the same things as the play process itself.

## API

### /upload/specified/path

Uploads the source code and makes it ready
to be used on the following path: `/specified/path`.
Note that your source code needs to have the
`Main(http.ResponseWriter, *http.Request)` function.

### /analyze/specified/path

If there is an error while building the source code,
it can be investigated on the following path:
`/analyze/specified/path`.

### /specified/path

After source code is successfully built, it handles requests
on the following path: `/specified/path`.

## How does it work?

Uploaded content is produced to a kafka topic.

Play is suscribed to consume that topic's messages.
Once a message is pulled, it builds the source code
provided by a message as a GO plugin and looks up
the Main function.

Any errors that occur while uploading are reported directly
to the client and also to play's logs (stdout) and metrics.

The error message createdd during this process is available
on the followingpath: `/analyze/specified/path`.
Make sure to always consult it.

If no error has occurred, you can a request to
the following: `/specified/path`.

Hurray! 🎉

### Replicas

In order to improve availability and to handle
more requests, multiple play servers can be
deployed.

The client usually sends requests to
a load balancer which forwards the request to
a replica chosen by some algorithm choosen.
E.g. Round robin or a more sophisticated
algorithm can be used. All replicas are the same so it
does not matter which one is picked for
the job.

#### Replicas can be added or removed

If more clients are expected, more play servers
can be added to handle more requests and vice
versa - play servers can be removed as there
are less requests to handle.

If a new replica is added, it starts pulling all
messages from the kafka topic - starting with the first
message and eventually will catch up with
the others.

### Result consistency

Note that play has an eventually consistent model.

The client only waits for a message to be sent (produced).
Client does not wait for a message to be processed,
for the source code to be built and so on. This happens
asynchronously. It means that result of requesting
`/specified/path` or `/analyze/specified/path`
right after uploading may not be available.

With multiple replicas and also depending on the deployed
load balancing strategy, the result might be
inconsistent - but eventually becomes consistent
as all replicas complete processing the message
and building.

## QUICKSTART

- Get kafka.
- Start the kafka environment.
- Create a topic to store uploaded content.
- Get play.
- Create `config.json` with the following content
  (to match your kafka broker and topic):

```json
{
  "kafkaBroker": "localhost:9092",
  "kafkaUploadTopic": "play-events",
  "apiKeys": [
    { "name": "main-api-key", "value": "ThatSecret1" },
    { "name": "user-api-key", "value": "ThatSecret2" }
  ]
}
```

```zsh
% play -addr=:8085 -config=config.json
```

- Create a GO file called `hello.go` with the following content:

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

- Use curl to upload `hello.go`:

```zsh
% curl --data-binary @hello.go -H "Authorization: Bearer ThatSecret1" http://localhost:8085/upload/say/hello
```

Note that `/upload/specified/path` is `/upload/say/hello`. `/say/hello` is going to be used as a path for calling that uploaded API endpoint.

```zsh
% curl -H "Authorization: Bearer ThatSecret2" http://localhost:8085/say/hello\?name=Gopher
Hello, Gopher!
```

## Package Documentations

- <https://pkg.go.dev/github.com/martindrlik/play/auth>
- <https://pkg.go.dev/github.com/martindrlik/play/config>
- <https://pkg.go.dev/github.com/martindrlik/play/her>
- <https://pkg.go.dev/github.com/martindrlik/play/id>
- <https://pkg.go.dev/github.com/martindrlik/play/kafka>
- <https://pkg.go.dev/github.com/martindrlik/play/limit>
- <https://pkg.go.dev/github.com/martindrlik/play/measure>
- <https://pkg.go.dev/github.com/martindrlik/play/metrics>
- <https://pkg.go.dev/github.com/martindrlik/play/plugin>
