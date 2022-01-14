## jwixtac

`jwixtac` is a repo template for bootstrapping applications with jwixel's preferred webstack.

### Setup

`boosc` has dependencies on [go](https://go.dev/), [nodejs](https://nodejs.org/en/), [GNU make](https://www.gnu.org/software/make/),
and [ffmpeg](https://www.ffmpeg.org/). Installing these dependencies is left as an exercise for the user.

#### Building

```go
make
```

#### Run dev UI

By running the Svelte dev target, hot-reloading is enabled in the http clients connected to the Go server. This is
highly recommended for development of the UI.

```sh
make run-dev
```

### About the stack

The stack uses go's standard library HTTP server with the gorilla [mux](https://github.com/gorilla/mux) and 
[websocket](https://github.com/gorilla/websocket) libraries.  The UI is a [Svelte](https://svelte.dev/) SPA.
