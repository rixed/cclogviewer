package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatToolResult(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "empty",
			content: "  \n",
			want:    "",
		},
		{
			name:    "plain text keeps whitespace",
			content: "line 1\n  indented <b>",
			want:    `<pre class="tool-result-text">line 1` + "\n" + `  indented &lt;b&gt;</pre>`,
		},
		{
			name:    "result wrapper is unwrapped",
			content: `{"result":"yadayada\nyada"}`,
			want:    `<pre class="tool-result-text">yadayada` + "\n" + `yada</pre>`,
		},
		{
			name:    "result wrapper around non-string is indented",
			content: `{"result":[1,2]}`,
			want:    `<pre class="tool-result-json">[` + "\n  1,\n  2\n]</pre>",
		},
		{
			name:    "other JSON is indented",
			content: `{"b":1,"a":"x"}`,
			want:    `<pre class="tool-result-json">{` + "\n  &#34;b&#34;: 1,\n  &#34;a&#34;: &#34;x&#34;\n}</pre>",
		},
		{
			name:    "bracketed notice",
			content: "[Request interrupted by user]",
			want:    `<span style="color: #999; font-style: italic;">Request interrupted by user</span>`,
		},
		{
			name:    "JSON array is not mistaken for a notice",
			content: `["a"]`,
			want:    `<pre class="tool-result-json">[` + "\n  &#34;a&#34;\n]</pre>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, string(formatToolResult(tt.content)))
		})
	}
}
