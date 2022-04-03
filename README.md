protoc-gen-typescript-http
==========================

Generates Typescript types and service clients from protobuf definitions annotated with [http rules](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto). The generated types follow the [canonical JSON encoding](https://developers.google.com/protocol-buffers/docs/proto3#json).

**Experimental**: This library is under active development and breaking changes to config files, APIs and generated code are expected between releases.

Using the plugin
----------------

For examples of correctly annotated protobuf defintions and the generated code, look at [examples](./examples).

### Install the plugin

```bash
go get go.einride.tech/protoc-gen-typescript-http
```

Or download a prebuilt binary from [releases](./releases).

### Invocation

```bash
protoc 
  --typescript-http_out [OUTPUT DIR] \
  [.proto files ...]
```

---

The generated clients can be used with any HTTP client that returns a Promise containing JSON data.

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
