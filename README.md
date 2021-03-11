# protoc-gen-typescript-http

Generates Typescript types and service clients from protobuf
definitions annotated with [http rules][httprule]. The generated
types follow the [canonical JSON encoding][jsonmapping].

[httprule]: https://github.com/googleapis/googleapis/blob/master/google/api/http.proto
[jsonmapping]: https://developers.google.com/protocol-buffers/docs/proto3#json

**Experimental**: This library is under active development and breaking 
changes to config files, APIs and generated code are expected between releases.

## Using the plugin

For examples of correctly annotated protobuf defintions and the 
generated code, look at [examples][examples].



### Install the plugin

```bash
go get github.com/einride/protoc-gen-typescript-http
```

Or download a prebuilt binary from [releases][releases].

### Invocation

```bash
protoc 
  --typescript-http_out [OUTPUT DIR] \
  [.proto files ...]
```

[examples]: ./examples
[releases]: ./releases

---

The generated clients can be used with any HTTP client that
returns a Promise containing JSON data.

```typescript
const rootUrl = "...";

type Request = {
  path: string,
  method: string,
  body: string | null
}

function fetchRequestHandler({path, method, body}: Request) {
  return fetch(rootUrl + path, {method, body}).then(response => response.json())
}

export function siteClient() {
  return createShipperServiceClient(fetchRequestHandler);
}
```
