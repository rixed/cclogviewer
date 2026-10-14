package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessUserMessageToolResult(t *testing.T) {
	tests := []struct {
		name    string
		content interface{}
		want    string
	}{
		{
			name:    "string content",
			content: "plain",
			want:    "plain",
		},
		{
			name: "all text blocks are joined",
			content: []interface{}{
				map[string]interface{}{"type": "text", "text": "first"},
				map[string]interface{}{"type": "image"},
				map[string]interface{}{"type": "text", "text": "second"},
			},
			want: "first\nsecond",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{"type": "tool_result", "content": tt.content},
				},
			}
			assert.Equal(t, tt.want, ProcessUserMessage(msg))
		})
	}
}
