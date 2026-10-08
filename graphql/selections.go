package graphql

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/OpenUdon/apitools/internal/sourceguard"
)

// Selection arguments need native typed values, rather than the lossy text used
// in prompt summaries. Sorting map encodings permits argument/object field order
// to differ while keeping lists, variables, enums and string values distinct.
func (p *graphQLParser) parseSelectionArguments() string {
	if !p.consume("(") {
		return ""
	}
	values := map[string]any{}
	for !p.done() && !p.peekValue(")") {
		name := p.consumeName()
		if name == "" || !p.consume(":") {
			p.err = fmt.Errorf("graphql: invalid selection argument")
			return ""
		}
		if _, duplicate := values[name]; duplicate {
			p.err = fmt.Errorf("graphql: duplicate selection argument")
			return ""
		}
		value, ok := p.selectionValue(0)
		if !ok {
			p.err = fmt.Errorf("graphql: invalid selection argument value")
			return ""
		}
		values[name] = value
	}
	if !p.consume(")") {
		p.err = fmt.Errorf("graphql: unterminated selection arguments")
		return ""
	}
	if len(values) == 0 {
		return ""
	}
	data, _ := json.Marshal(values)
	return string(data)
}

func (p *graphQLParser) selectionValue(depth int) (any, bool) {
	if p.done() || depth >= sourceguard.MaxNestingDepth {
		return nil, false
	}
	token := p.advance()
	if token.kind == tokenString {
		if !utf8.ValidString(token.raw) {
			return nil, false
		}
		if strings.HasPrefix(token.raw, `"""`) {
			return []any{"string", selectionBlockStringValue(token.raw)}, true
		}
		value, ok := selectionQuotedStringValue(token.raw)
		return []any{"string", value}, ok
	}
	if token.kind == tokenName {
		return []any{"literal", token.value}, true
	}
	if token.kind != tokenPunct {
		return nil, false
	}
	switch token.value {
	case "$":
		name := p.consumeName()
		return []any{"variable", name}, name != ""
	case "[":
		values := []any{}
		for !p.done() && !p.peekValue("]") {
			value, ok := p.selectionValue(depth + 1)
			if !ok {
				return nil, false
			}
			values = append(values, value)
		}
		return []any{"list", values}, p.consume("]")
	case "{":
		values := map[string]any{}
		for !p.done() && !p.peekValue("}") {
			name := p.consumeName()
			if name == "" || !p.consume(":") {
				return nil, false
			}
			if _, duplicate := values[name]; duplicate {
				return nil, false
			}
			value, ok := p.selectionValue(depth + 1)
			if !ok {
				return nil, false
			}
			values[name] = value
		}
		return []any{"object", values}, p.consume("}")
	}
	return nil, false
}

// GraphQL block strings remove common indentation after the first line and
// leading/trailing blank lines. Their semantic value can equal a quoted string.
func selectionBlockStringValue(raw string) string {
	value := strings.ReplaceAll(raw[3:len(raw)-3], `\"""`, `"""`)
	value = strings.ReplaceAll(strings.ReplaceAll(value, "\r\n", "\n"), "\r", "\n")
	lines := strings.Split(value, "\n")
	indent := -1
	for _, line := range lines[1:] {
		count := len(line) - len(strings.TrimLeft(line, " \t"))
		if count != len(line) && (indent < 0 || count < indent) {
			indent = count
		}
	}
	if indent > 0 {
		for i := 1; i < len(lines); i++ {
			lines[i] = lines[i][min(indent, len(lines[i])):]
		}
	}
	for len(lines) > 0 && strings.Trim(lines[0], " \t") == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.Trim(lines[len(lines)-1], " \t") == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}

// Native GraphQL escapes include variable-width Unicode, and reject lone
// surrogates. JSON's replacement behavior cannot establish argument equality.
func selectionQuotedStringValue(raw string) (string, bool) {
	if len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' {
		return "", false
	}
	var out strings.Builder
	for pos := 1; pos < len(raw)-1; {
		if raw[pos] != '\\' {
			value, size := utf8.DecodeRuneInString(raw[pos : len(raw)-1])
			if value == utf8.RuneError && size == 1 || value == '\n' || value == '\r' || value == '"' {
				return "", false
			}
			out.WriteRune(value)
			pos += size
			continue
		}
		pos++
		if pos >= len(raw)-1 {
			return "", false
		}
		escape := raw[pos]
		pos++
		switch escape {
		case '"', '\\', '/':
			out.WriteByte(escape)
		case 'b':
			out.WriteByte('\b')
		case 'f':
			out.WriteByte('\f')
		case 'n':
			out.WriteByte('\n')
		case 'r':
			out.WriteByte('\r')
		case 't':
			out.WriteByte('\t')
		case 'u':
			braced := pos < len(raw)-1 && raw[pos] == '{'
			if braced {
				pos++
			}
			start := pos
			var value rune
			for pos < len(raw)-1 && (braced && raw[pos] != '}' || !braced && pos-start < 4) {
				hex := raw[pos]
				var digit rune
				switch {
				case hex >= '0' && hex <= '9':
					digit = rune(hex - '0')
				case hex >= 'a' && hex <= 'f':
					digit = rune(hex - 'a' + 10)
				case hex >= 'A' && hex <= 'F':
					digit = rune(hex - 'A' + 10)
				default:
					return "", false
				}
				value = value*16 + digit
				if value > utf8.MaxRune {
					return "", false
				}
				pos++
			}
			if pos == start || !braced && pos-start != 4 {
				return "", false
			}
			if braced {
				if pos >= len(raw)-1 || raw[pos] != '}' {
					return "", false
				}
				pos++
			}
			if value >= 0xD800 && value <= 0xDBFF && !braced {
				if pos+6 > len(raw)-1 || raw[pos:pos+2] != `\u` {
					return "", false
				}
				low, err := strconv.ParseUint(raw[pos+2:pos+6], 16, 16)
				if err != nil || low < 0xDC00 || low > 0xDFFF {
					return "", false
				}
				value = utf16.DecodeRune(value, rune(low))
				pos += 6
			} else if value >= 0xD800 && value <= 0xDFFF {
				return "", false
			}
			out.WriteRune(value)
		default:
			return "", false
		}
	}
	return out.String(), true
}
