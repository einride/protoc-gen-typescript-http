package httprule

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_ParseTemplate(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		input string
		path  Template
	}{
		{
			input: "/v1/messages",
			path: Template{
				Segments: []Segment{
					{Kind: SegmentKindLiteral, Literal: "v1"},
					{Kind: SegmentKindLiteral, Literal: "messages"},
				},
			},
		},
		{
			input: "/v1/messages:peek",
			path: Template{
				Segments: []Segment{
					{Kind: SegmentKindLiteral, Literal: "v1"},
					{Kind: SegmentKindLiteral, Literal: "messages"},
				},
				Verb: "peek",
			},
		},
		{
			input: "/{id}",
			path: Template{
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
			path: Template{
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
			path: Template{
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
			path: Template{
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
			path: Template{
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
			path: Template{
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
	for _, tt := range []struct {
		template string
		expected string
	}{
		{template: "", expected: "expected token '/' at position 0, found EOF"},
		{template: "//", expected: "expected literal at position 1, found '/'"},
		{template: "/v1:", expected: "expected literal at position 3, found EOF"},
		{template: "/v1/:", expected: "expected literal at position 4, found ':'"},
		{template: "/{name=messages/{id}}", expected: "nested variable segment is not allowed"},
		{template: "/**/*", expected: "'**' only allowed as last part of template"},
		{template: "/v1/messages/*", expected: "'*' must only be used in variables"},
		{template: "/v1/{id}/{id}", expected: "variable 'id' bound multiple times"},
	} {
		t.Run(tt.template, func(t *testing.T) {
			t.Parallel()
			_, err := ParseTemplate(tt.template)
			assert.Check(t, err != nil)
			assert.ErrorContains(t, err, tt.expected)
		})
	}
}
