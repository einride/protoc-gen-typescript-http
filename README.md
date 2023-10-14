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
  --typescript-http_opt use_enum_numbers=true,use_multi_line_comment=true
  [.proto files ...]
```

#### Support options

```ts
// UseProtoNames controls the casing of generated field names.
// If set to true, fields will use proto names (typically snake_case).
// If omitted or set to false, fields will use JSON names (typically camelCase).
use_proto_names: bool

// UseEnumNumbers emits enum values as numbers.
use_enum_numbers: bool

// The method names of service methods naming case.
// Only work when `UseEnumNumbers=true`
// opt:
// camelcase: convert name to lower camel case like `camelCase`
// pascalcase: convert name to pascalcase like `PascalCase`
// default is pascalcase
enum_field_naming: string

// Generate comments as multiline comments.
// multiline comments: /** ... */
// single line comments: // ...
use_multi_line_comment: bool

// force add `undefined` to message field.
// default true
force_message_field_undefinable: bool

// If set to true, body will be JSON.stringify before send
// default true
use_body_stringify: bool

// The method names of service methods naming case.
// opt:
// camelcase: convert name to lower camel case like `camelCase`
// pascalcase: convert name to pascalcase like `PascalCase`
// default is pascalcase
service_method_naming: 'camelcase' | 'pascalcase'

// If set to true, field int64 and uint64 will convert to string
force_long_as_string: bool
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

// This is optional
type RequestOptions = {
  useCache?: boolean;
}

function fetchRequestHandler({path, method, body}: Request & RequestOptions) {
  return fetch(rootUrl + path, {method, body}).then(response => response.json())
}

export function siteClient() {
  // This Generics is optional
  return createShipperServiceClient<RequestOptions>(fetchRequestHandler);
}
```
