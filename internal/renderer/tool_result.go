package renderer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strings"
)

var ansiSGRPattern = regexp.MustCompile(`\x1b\[\d+m`)

// formatNotice renders bracketed notices such as "[Request interrupted by user]"
// in a muted style. ok is false when content is not such a notice.
func formatNotice(content string) (out template.HTML, ok bool) {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") ||
		ansiSGRPattern.MatchString(content) {
		return "", false
	}
	stripped := trimmed[1 : len(trimmed)-1]
	return template.HTML(fmt.Sprintf(`<span style="color: #999; font-style: italic;">%s</span>`,
		html.EscapeString(stripped))), true
}

// formatToolResult renders a tool result as preformatted text. JSON results are
// pretty-printed, and a lone {"result": ...} wrapper (as emitted by some MCP
// servers) is unwrapped so that string results read as plain text.
func formatToolResult(content string) template.HTML {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ""
	}

	if json.Valid([]byte(trimmed)) {
		raw := json.RawMessage(trimmed)
		var wrapper map[string]json.RawMessage
		if json.Unmarshal(raw, &wrapper) == nil && len(wrapper) == 1 && wrapper["result"] != nil {
			raw = wrapper["result"]
		}
		var s string
		if json.Unmarshal(raw, &s) == nil {
			return formatToolResultText(s)
		}
		var indented bytes.Buffer
		if json.Indent(&indented, raw, "", "  ") == nil {
			return template.HTML(`<pre class="tool-result-json">` +
				html.EscapeString(indented.String()) + `</pre>`)
		}
	}

	if notice, ok := formatNotice(content); ok {
		return notice
	}
	return formatToolResultText(content)
}

func formatToolResultText(text string) template.HTML {
	return template.HTML(`<pre class="tool-result-text">` +
		ConvertANSIToHTML(strings.TrimRight(text, "\n")) + `</pre>`)
}
