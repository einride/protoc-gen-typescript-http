// Message
export type Message = {
	// double
	double?: number;
	// float
	float?: number;
	// int32
	int32?: number;
	// int64
	int64?: number;
	// uint32
	uint32?: number;
	// uint64
	uint64?: number;
	// sint32
	sint32?: unknown;
	// sint64
	sint64?: unknown;
	// fixed32
	fixed32?: number;
	// fixed64
	fixed64?: number;
	// sfixed32
	sfixed32?: number;
	// sfixed64
	sfixed64?: number;
	// bool
	bool?: boolean;
	// string
	string?: string;
	// bytes
	bytes?: string;
	// enum
	enum?: unknown;
	// message
	message?: unknown;
	// repeated_double
	repeatedDouble?: number[];
	// repeated_float
	repeatedFloat?: number[];
	// repeated_int32
	repeatedInt32?: number[];
	// repeated_int64
	repeatedInt64?: number[];
	// repeated_uint32
	repeatedUint32?: number[];
	// repeated_uint64
	repeatedUint64?: number[];
	// repeated_sint32
	repeatedSint32?: unknown[];
	// repeated_sint64
	repeatedSint64?: unknown[];
	// repeated_fixed32
	repeatedFixed32?: number[];
	// repeated_fixed64
	repeatedFixed64?: number[];
	// repeated_sfixed32
	repeatedSfixed32?: number[];
	// repeated_sfixed64
	repeatedSfixed64?: number[];
	// repeated_bool
	repeatedBool?: boolean[];
	// repeated_string
	repeatedString?: string[];
	// repeated_bytes
	repeatedBytes?: string[];
	// repeated_enum
	repeatedEnum?: unknown[];
	// repeated_message
	repeatedMessage?: unknown[];
	// map_string_string
	mapStringString?: { [key: string]: string};
	// map_string_message
	mapStringMessage?: { [key: string]: unknown};
	// oneof_string
	oneofString?: string;
	// oneof_enum
	oneofEnum?: unknown;
	// oneof_message1
	oneofMessage1?: unknown;
	// oneof_message2
	oneofMessage2?: unknown;
};

// NestedMessage
export type Message_NestedMessage = {
	// nested_message.string
	string?: string;
};

// NestedEnum
export type Message_NestedEnum = 
	// NESTEDENUM_UNSPECIFIED
	| "NESTEDENUM_UNSPECIFIED"

// Enum
export type Enum = 
	// ENUM_UNSPECIFIED
	| "ENUM_UNSPECIFIED"
	// ENUM_ONE
	| "ENUM_ONE"
	// ENUM_TWO
	| "ENUM_TWO"

