// Code generated by protoc-gen-typescript-http. DO NOT EDIT.
/* eslint-disable camelcase */

// Message
export type Message = {
  forwardedMessage: einrideexamplesyntaxv1_Message | undefined;
  forwardedEnum: einrideexamplesyntaxv1_Enum | undefined;
};

// Message
export type einrideexamplesyntaxv1_Message = {
  // double
  double: number | undefined;
  // float
  float: number | undefined;
  // int32
  int32: number | undefined;
  // int64
  int64: number | undefined;
  // uint32
  uint32: number | undefined;
  // uint64
  uint64: number | undefined;
  // sint32
  sint32: number | undefined;
  // sint64
  sint64: number | undefined;
  // fixed32
  fixed32: number | undefined;
  // fixed64
  fixed64: number | undefined;
  // sfixed32
  sfixed32: number | undefined;
  // sfixed64
  sfixed64: number | undefined;
  // bool
  bool: boolean | undefined;
  // string
  string: string | undefined;
  // bytes
  bytes: string | undefined;
  // enum
  enum: einrideexamplesyntaxv1_Enum | undefined;
  // message
  message: einrideexamplesyntaxv1_Message | undefined;
  // optional double
  optionalDouble?: number;
  // optional float
  optionalFloat?: number;
  // optional int32
  optionalInt32?: number;
  // optional int64
  optionalInt64?: number;
  // optional uint32
  optionalUint32?: number;
  // optional uint64
  optionalUint64?: number;
  // optional sint32
  optionalSint32?: number;
  // optional sint64
  optionalSint64?: number;
  // optional fixed32
  optionalFixed32?: number;
  // optional fixed64
  optionalFixed64?: number;
  // optional sfixed32
  optionalSfixed32?: number;
  // optional sfixed64
  optionalSfixed64?: number;
  // optional bool
  optionalBool?: boolean;
  // optional string
  optionalString?: string;
  // optional bytes
  optionalBytes?: string;
  // optional enum
  optionalEnum?: einrideexamplesyntaxv1_Enum;
  // optional message
  optionalMessage?: einrideexamplesyntaxv1_Message;
  // repeated_double
  repeatedDouble: number[] | undefined;
  // repeated_float
  repeatedFloat: number[] | undefined;
  // repeated_int32
  repeatedInt32: number[] | undefined;
  // repeated_int64
  repeatedInt64: number[] | undefined;
  // repeated_uint32
  repeatedUint32: number[] | undefined;
  // repeated_uint64
  repeatedUint64: number[] | undefined;
  // repeated_sint32
  repeatedSint32: number[] | undefined;
  // repeated_sint64
  repeatedSint64: number[] | undefined;
  // repeated_fixed32
  repeatedFixed32: number[] | undefined;
  // repeated_fixed64
  repeatedFixed64: number[] | undefined;
  // repeated_sfixed32
  repeatedSfixed32: number[] | undefined;
  // repeated_sfixed64
  repeatedSfixed64: number[] | undefined;
  // repeated_bool
  repeatedBool: boolean[] | undefined;
  // repeated_string
  repeatedString: string[] | undefined;
  // repeated_bytes
  repeatedBytes: string[] | undefined;
  // repeated_enum
  repeatedEnum: einrideexamplesyntaxv1_Enum[] | undefined;
  // repeated_message
  repeatedMessage: einrideexamplesyntaxv1_Message[] | undefined;
  // map_string_string
  mapStringString: { [key: string]: string } | undefined;
  // map_string_message
  mapStringMessage: { [key: string]: einrideexamplesyntaxv1_Message } | undefined;
  // oneof_string
  oneofString?: string;
  // oneof_enum
  oneofEnum?: einrideexamplesyntaxv1_Enum;
  // oneof_message1
  oneofMessage1?: einrideexamplesyntaxv1_Message;
  // oneof_message2
  oneofMessage2?: einrideexamplesyntaxv1_Message;
  // any
  any: wellKnownAny | undefined;
  // repeated_any
  repeatedAny: wellKnownAny[] | undefined;
  // duration
  duration: wellKnownDuration | undefined;
  // repeated_duration
  repeatedDuration: wellKnownDuration[] | undefined;
  // empty
  empty: wellKnownEmpty | undefined;
  // repeated_empty
  repeatedEmpty: wellKnownEmpty[] | undefined;
  // field_mask
  fieldMask: wellKnownFieldMask | undefined;
  // repeated_field_mask
  repeatedFieldMask: wellKnownFieldMask[] | undefined;
  // struct
  struct: wellKnownStruct | undefined;
  // repeated_struct
  repeatedStruct: wellKnownStruct[] | undefined;
  // value
  value: wellKnownValue | undefined;
  // repeated_value
  repeatedValue: wellKnownValue[] | undefined;
  // null_value
  nullValue: wellKnownNullValue | undefined;
  // repeated_null_value
  repeatedNullValue: wellKnownNullValue[] | undefined;
  // list_value
  listValue: wellKnownListValue | undefined;
  // repeated_list_value
  repeatedListValue: wellKnownListValue[] | undefined;
  // bool_value
  boolValue: wellKnownBoolValue | undefined;
  // repeated_bool_value
  repeatedBoolValue: wellKnownBoolValue[] | undefined;
  // bytes_value
  bytesValue: wellKnownBytesValue | undefined;
  // repeated_bytes_value
  repeatedBytesValue: wellKnownBytesValue[] | undefined;
  // double_value
  doubleValue: wellKnownDoubleValue | undefined;
  // repeated_double_value
  repeatedDoubleValue: wellKnownDoubleValue[] | undefined;
  // float_value
  floatValue: wellKnownFloatValue | undefined;
  // repeated_float_value
  repeatedFloatValue: wellKnownFloatValue[] | undefined;
  // int32_value
  int32Value: wellKnownInt32Value | undefined;
  // repeated_int32_value
  repeatedInt32Value: wellKnownInt32Value[] | undefined;
  // int64_value
  int64Value: wellKnownInt64Value | undefined;
  // repeated_int64_value
  repeatedInt64Value: wellKnownInt64Value[] | undefined;
  // uint32_value
  uint32Value: wellKnownUInt32Value | undefined;
  // repeated_uint32_value
  repeatedUint32Value: wellKnownUInt32Value[] | undefined;
  // uint64_value
  uint64Value: wellKnownUInt64Value | undefined;
  // repeated_uint64_value
  repeatedUint64Value: wellKnownUInt64Value[] | undefined;
  // string_value
  stringValue: wellKnownUInt64Value | undefined;
  // repeated_string_value
  repeatedStringValue: wellKnownStringValue[] | undefined;
};

// Enum
export type einrideexamplesyntaxv1_Enum =
  // ENUM_UNSPECIFIED
  | "ENUM_UNSPECIFIED"
  // ENUM_ONE
  | "ENUM_ONE"
  // ENUM_TWO
  | "ENUM_TWO";
// If the Any contains a value that has a special JSON mapping,
// it will be converted as follows:
// {"@type": xxx, "value": yyy}.
// Otherwise, the value will be converted into a JSON object,
// and the "@type" field will be inserted to indicate the actual data type.
interface wellKnownAny {
  "@type": string;
  [key: string]: unknown;
}

// Generated output always contains 0, 3, 6, or 9 fractional digits,
// depending on required precision, followed by the suffix "s".
// Accepted are any fractional digits (also none) as long as they fit
// into nano-seconds precision and the suffix "s" is required.
type wellKnownDuration = string;

// An empty JSON object
type wellKnownEmpty = Record<never, never>;

// In JSON, a field mask is encoded as a single string where paths are
// separated by a comma. Fields name in each path are converted
// to/from lower-camel naming conventions.
// As an example, consider the following message declarations:
//
//     message Profile {
//       User user = 1;
//       Photo photo = 2;
//     }
//     message User {
//       string display_name = 1;
//       string address = 2;
//     }
//
// In proto a field mask for `Profile` may look as such:
//
//     mask {
//       paths: "user.display_name"
//       paths: "photo"
//     }
//
// In JSON, the same mask is represented as below:
//
//     {
//       mask: "user.displayName,photo"
//     }
type wellKnownFieldMask = string;

// Any JSON value.
type wellKnownStruct = Record<string, unknown>;

type wellKnownValue = unknown;

type wellKnownNullValue = null;

type wellKnownListValue = wellKnownValue[];

type wellKnownBoolValue = boolean | null;

type wellKnownBytesValue = string | null;

type wellKnownDoubleValue = number | null;

type wellKnownFloatValue = number | null;

type wellKnownInt32Value = number | null;

type wellKnownInt64Value = number | null;

type wellKnownUInt32Value = number | null;

type wellKnownUInt64Value = number | null;

type wellKnownStringValue = string | null;

// NestedMessage
export type einrideexamplesyntaxv1_Message_NestedMessage = {
  // nested_message.string
  string: string | undefined;
};

// NestedEnum
export type einrideexamplesyntaxv1_Message_NestedEnum =
  // NESTEDENUM_UNSPECIFIED
  "NESTEDENUM_UNSPECIFIED";

// @@protoc_insertion_point(typescript-http-eof)