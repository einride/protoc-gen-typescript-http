package httprule

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_ParseTemplate(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		input string
		path  *Template
	}{
		{
			input: "/*",
			path: &Template{
				Segments: []Segment{
					{Kind: SegmentKindMatchSingle},
				},
			},
		},
		{
			input: "/**",
			path: &Template{
				Segments: []Segment{
					{Kind: SegmentKindMatchMultiple},
				},
			},
		},
		{
			input: "/*/*",
			path: &Template{
				Segments: []Segment{
					{Kind: SegmentKindMatchSingle},
					{Kind: SegmentKindMatchSingle},
				},
			},
		},
		{
			input: "/**:peek",
			path: &Template{
				Segments: []Segment{
					{Kind: SegmentKindMatchMultiple},
				},
				Verb: "peek",
			},
		},
		{
			input: "/v1/messages",
			path: &Template{
				Segments: []Segment{
					{Kind: SegmentKindLiteral, Literal: "v1"},
					{Kind: SegmentKindLiteral, Literal: "messages"},
				},
			},
		},
		{
			input: "/v1/**",
			path: &Template{
				Segments: []Segment{
					{Kind: SegmentKindLiteral, Literal: "v1"},
					{Kind: SegmentKindMatchMultiple},
				},
			},
		},
		{
			input: "/v1/**:peek",
			path: &Template{
				Segments: []Segment{
					{Kind: SegmentKindLiteral, Literal: "v1"},
					{Kind: SegmentKindMatchMultiple},
				},
				Verb: "peek",
			},
		},
		{
			input: "/{id}",
			path: &Template{
				Segments: []Segment{
					{
						Kind: SegmentKindVariable,
						Variable: VariableSegment{
							FieldPath: []string{"id"},
							Segments: []Segment{
								{Kind: SegmentKindMatchSingle},
							},
						},
					},
				},
			},
		},
		{
			input: "/{message.id}",
			path: &Template{
				Segments: []Segment{
					{
						Kind: SegmentKindVariable,
						Variable: VariableSegment{
							FieldPath: []string{"message", "id"},
							Segments: []Segment{
								{Kind: SegmentKindMatchSingle},
							},
						},
					},
				},
			},
		},
		{
			input: "/{id=messages/*}",
			path: &Template{
				Segments: []Segment{
					{
						Kind: SegmentKindVariable,
						Variable: VariableSegment{
							FieldPath: []string{"id"},
							Segments: []Segment{
								{Kind: SegmentKindLiteral, Literal: "messages"},
								{Kind: SegmentKindMatchSingle},
							},
						},
					},
				},
			},
		},
		{
			input: "/{id=messages/*/threads/*}",
			path: &Template{
				Segments: []Segment{
					{
						Kind: SegmentKindVariable,
						Variable: VariableSegment{
							FieldPath: []string{"id"},
							Segments: []Segment{
								{Kind: SegmentKindLiteral, Literal: "messages"},
								{Kind: SegmentKindMatchSingle},
								{Kind: SegmentKindLiteral, Literal: "threads"},
								{Kind: SegmentKindMatchSingle},
							},
						},
					},
				},
			},
		},
		{
			input: "/{id=**}",
			path: &Template{
				Segments: []Segment{
					{
						Kind: SegmentKindVariable,
						Variable: VariableSegment{
							FieldPath: []string{"id"},
							Segments: []Segment{
								{Kind: SegmentKindMatchMultiple},
							},
						},
					},
				},
			},
		},
		{
			input: "/v1/messages/{message}/threads/{thread}",
			path: &Template{
				Segments: []Segment{
					{Kind: SegmentKindLiteral, Literal: "v1"},
					{Kind: SegmentKindLiteral, Literal: "messages"},
					{
						Kind: SegmentKindVariable,
						Variable: VariableSegment{
							FieldPath: []string{"message"},
							Segments: []Segment{
								{Kind: SegmentKindMatchSingle},
							},
						},
					},
					{Kind: SegmentKindLiteral, Literal: "threads"},
					{
						Kind: SegmentKindVariable,
						Variable: VariableSegment{
							FieldPath: []string{"thread"},
							Segments: []Segment{
								{Kind: SegmentKindMatchSingle},
							},
						},
					},
				},
			},
		},
	} {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParseTemplate(tt.input)
			assert.NilError(t, err)
			assert.DeepEqual(t, tt.path, got)
		})
	}
}

func Test_ParseTemplate_Invalid(t *testing.T) {
	t.Parallel()
	for _, tt := range []string{
		"",
		"//",
		"/v1:",
		"/v1/:",
		"/{name=messages/{id}}",
		"/**/*",
		"v1/messages/*",
		"v1/{id}/{id}",
	} {
		tt := tt
		t.Run(tt, func(t *testing.T) {
			t.Parallel()
			_, err := ParseTemplate(tt)
			assert.Check(t, err != nil)
		})
	}
}
