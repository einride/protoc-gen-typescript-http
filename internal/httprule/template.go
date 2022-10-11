package httprule

import "fmt"

// Template represents a http path template.
//
// Example: `/v1/{name=books/*}:publish`.
type Template struct {
	Segments []Segment
	Verb     string
}

// Segment represents a single segment of a Template.
type Segment struct {
	Kind     SegmentKind
	Literal  string
	Variable VariableSegment
}

type SegmentKind int

const (
	SegmentKindLiteral SegmentKind = iota
	SegmentKindMatchSingle
	SegmentKindMatchMultiple
	SegmentKindVariable
)

// VariableSegment represents a variable segment.
type VariableSegment struct {
	FieldPath FieldPath
	Segments  []Segment
}

func ParseTemplate(s string) (Template, error) {
	p := &parser{
		content: s,
	}
	template, err := p.parse()
	if err != nil {
		return Template{}, err
	}
	if err := validate(template); err != nil {
		return Template{}, err
	}
	return template, nil
}

type parser struct {
	content string

	// The next pos in content to read
	pos int
	// The currently read rune in content
	tok rune
}

func (p *parser) parse() (Template, error) {
	// Grammar.
	// Template = "/" Segments [ Verb ] ;
	// Segments = Segment { "/" Segment } ;
	// Segment  = "*" | "**" | LITERAL | Variable ;
	// Variable = "{" FieldPath [ "=" Segments ] "}" ;
	// FieldPath = IDENT { "." IDENT } ;
	// Verb     = ":" LITERAL ;.
	p.next()
	if err := p.expect('/'); err != nil {
		return Template{}, err
	}
	segments, err := p.parseSegments()
	if err != nil {
		return Template{}, err
	}
	var verb string
	if p.tok == ':' {
		v, err := p.parseVerb()
		if err != nil {
			return Template{}, err
		}
		verb = v
	}
	if p.tok != -1 {
		return Template{}, fmt.Errorf("expected EOF, got %q", p.tok)
	}
	return Template{
		Segments: segments,
		Verb:     verb,
	}, nil
}

func (p *parser) parseSegments() ([]Segment, error) {
	seg, err := p.parseSegment()
	if err != nil {
		return nil, err
	}
	if p.tok == '/' {
		p.next()
		rest, err := p.parseSegments()
		if err != nil {
			return nil, err
		}
		return append([]Segment{seg}, rest...), nil
	}
	return []Segment{seg}, nil
}

func (p *parser) parseSegment() (Segment, error) {
	switch {
	case p.tok == '*' && p.peek() == '*':
		return p.parseMatchMultipleSegment(), nil
	case p.tok == '*':
		return p.parseMatchSingleSegment(), nil
	case p.tok == '{':
		return p.parseVariableSegment()
	default:
		return p.parseLiteralSegment()
	}
}

func (p *parser) parseMatchMultipleSegment() Segment {
	p.next()
	p.next()
	return Segment{
		Kind: SegmentKindMatchMultiple,
	}
}

func (p *parser) parseMatchSingleSegment() Segment {
	p.next()
	return Segment{
		Kind: SegmentKindMatchSingle,
	}
}

func (p *parser) parseLiteralSegment() (Segment, error) {
	lit, err := p.parseLiteral()
	if err != nil {
		return Segment{}, err
	}
	return Segment{
		Kind:    SegmentKindLiteral,
		Literal: lit,
	}, nil
}

func (p *parser) parseVariableSegment() (Segment, error) {
	if err := p.expect('{'); err != nil {
		return Segment{}, err
	}
	fieldPath, err := p.parseFieldPath()
	if err != nil {
		return Segment{}, err
	}
	segments := []Segment{
		{Kind: SegmentKindMatchSingle},
	}
	if p.tok == '=' {
		p.next()
		s, err := p.parseSegments()
		if err != nil {
			return Segment{}, err
		}
		segments = s
	}
	if err := p.expect('}'); err != nil {
		return Segment{}, err
	}
	return Segment{
		Kind: SegmentKindVariable,
		Variable: VariableSegment{
			FieldPath: fieldPath,
			Segments:  segments,
		},
	}, nil
}

func (p *parser) parseVerb() (string, error) {
	if err := p.expect(':'); err != nil {
		return "", err
	}
	return p.parseLiteral()
}

func (p *parser) parseFieldPath() ([]string, error) {
	fp, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if p.tok == '.' {
		p.next()
		rest, err := p.parseFieldPath()
		if err != nil {
			return nil, err
		}
		return append([]string{fp}, rest...), nil
	}
	return []string{fp}, nil
}

// parseLiteral consumes input as long as next token(s) belongs to pchars, as defined in RFC3986.
// Returns an error if not literal is found.
//
// https://www.ietf.org/rfc/rfc3986.txt, P.49
//
//	pchar         = unreserved / pct-encoded / sub-delims / ":" / "@"
//	unreserved    = ALPHA / DIGIT / "-" / "." / "_" / "~"
//	sub-delims    = "!" / "$" / "&" / "'" / "(" / ")"
//	              / "*" / "+" / "," / ";" / "="
//	pct-encoded   = "%" HEXDIG HEXDIG
func (p *parser) parseLiteral() (string, error) {
	var literal []rune
	startPos := p.pos
	for {
		if isSingleCharPChar(p.tok) {
			literal = append(literal, p.tok)
			p.next()
			continue
		}
		if p.tok == '%' && isHexDigit(p.peekN(1)) && isHexDigit(p.peekN(2)) {
			literal = append(literal, p.tok)
			p.next()
			literal = append(literal, p.tok)
			p.next()
			literal = append(literal, p.tok)
			p.next()
			continue
		}
		break
	}
	if len(literal) == 0 {
		return "", fmt.Errorf("expected literal at position %d, found %s", startPos-1, p.tokenString())
	}
	return string(literal), nil
}

func (p *parser) parseIdent() (string, error) {
	var ident []rune
	startPos := p.pos
	for {
		if isAlpha(p.tok) || isDigit(p.tok) || p.tok == '_' {
			ident = append(ident, p.tok)
			p.next()
			continue
		}
		break
	}
	if len(ident) == 0 {
		return "", fmt.Errorf("expected identifier at position %d, found %s", startPos-1, p.tokenString())
	}
	return string(ident), nil
}

func (p *parser) next() {
	if p.pos < len(p.content) {
		p.tok = rune(p.content[p.pos])
		p.pos++
	} else {
		p.tok = -1
		p.pos = len(p.content)
	}
}

func (p parser) tokenString() string {
	if p.tok == -1 {
		return "EOF"
	}
	return fmt.Sprintf("%q", p.tok)
}

func (p *parser) peek() rune {
	return p.peekN(1)
}

func (p *parser) peekN(n int) rune {
	if offset := p.pos + n - 1; offset < len(p.content) {
		return rune(p.content[offset])
	}
	return -1
}

func (p *parser) expect(r rune) error {
	if p.tok != r {
		return fmt.Errorf("expected token %q at position %d, found %s", r, p.pos, p.tokenString())
	}
	p.next()
	return nil
}

// https://www.ietf.org/rfc/rfc3986.txt, P.49
//
//	pchar         = unreserved / pct-encoded / sub-delims / ":" / "@"
//	unreserved    = ALPHA / DIGIT / "-" / "." / "_" / "~"
//	sub-delims    = "!" / "$" / "&" / "'" / "(" / ")"
//	              / "*" / "+" / "," / ";" / "="
//	pct-encoded   = "%" HEXDIG HEXDIG
func isSingleCharPChar(r rune) bool {
	if isAlpha(r) || isDigit(r) {
		return true
	}
	switch r {
	case '@', '-', '.', '_', '~', '!',
		'$', '&', '\'', '(', ')', '*', '+',
		',', ';', '=': // ':'
		return true
	}
	return false
}

func isAlpha(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isHexDigit(r rune) bool {
	switch {
	case '0' <= r && r <= '9':
		return true
	case 'A' <= r && r <= 'F':
		return true
	case 'a' <= r && r <= 'f':
		return true
	}
	return false
}

// validate validates parts of the template that are
// allowed by the grammar, but disallowed in practice.
//
// - nested variable segments
// - '**' for segments other than the last.
func validate(t Template) error {
	// check for nested variable segments
	for _, s1 := range t.Segments {
		if s1.Kind != SegmentKindVariable {
			continue
		}
		for _, s2 := range s1.Variable.Segments {
			if s2.Kind == SegmentKindVariable {
				return fmt.Errorf("nested variable segment is not allowed")
			}
		}
	}
	// check for '**' that are not the last part of the template
	for i, s := range t.Segments {
		if i == len(t.Segments)-1 {
			continue
		}
		if s.Kind == SegmentKindMatchMultiple {
			return fmt.Errorf("'**' only allowed as last part of template")
		}
		if s.Kind == SegmentKindVariable {
			for _, s2 := range s.Variable.Segments {
				if s2.Kind == SegmentKindMatchMultiple {
					return fmt.Errorf("'**' only allowed as last part of template")
				}
			}
		}
	}
	// check for variable where '**' is not last part
	for _, s := range t.Segments {
		if s.Kind != SegmentKindVariable {
			continue
		}
		for i, s2 := range s.Variable.Segments {
			if i == len(s.Variable.Segments)-1 {
				continue
			}
			if s2.Kind == SegmentKindMatchMultiple {
				return fmt.Errorf("'**' only allowed as the last part of the template")
			}
		}
	}
	// check for top level expansions
	for _, s := range t.Segments {
		if s.Kind == SegmentKindMatchSingle {
			return fmt.Errorf("'*' must only be used in variables")
		}
		if s.Kind == SegmentKindMatchMultiple {
			return fmt.Errorf("'**' must only be used in variables")
		}
	}
	// check for duplicate variable bindings
	seen := make(map[string]struct{})
	for _, s := range t.Segments {
		if s.Kind == SegmentKindVariable {
			field := s.Variable.FieldPath.String()
			if _, ok := seen[s.Variable.FieldPath.String()]; ok {
				return fmt.Errorf("variable '%s' bound multiple times", field)
			}
			seen[field] = struct{}{}
		}
	}
	return nil
}
