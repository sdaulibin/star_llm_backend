package dify

import (
	"bufio"
	"io"
	"strings"
)

// SSEEvent 表示服务器发送的事件
type SSEEvent struct {
	Event string
	Data  string
}

// SSEReader 是一个服务器发送事件的读取器
type SSEReader struct {
	reader *bufio.Reader
}

// NewSSEReader 创建一个新的 SSE 读取器
func NewSSEReader(r io.Reader) *SSEReader {
	return &SSEReader{
		reader: bufio.NewReader(r),
	}
}

// ReadEvent 读取下一个 SSE 事件
func (r *SSEReader) ReadEvent() (*SSEEvent, error) {
	event := &SSEEvent{}
	var data strings.Builder

	for {
		line, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\n\r")
		if line == "" {
			// 空行表示事件结束
			if data.Len() > 0 {
				event.Data = data.String()
				return event, nil
			}
			continue
		}

		if strings.HasPrefix(line, "event:") {
			event.Event = strings.TrimSpace(line[6:])
		} else if strings.HasPrefix(line, "data:") {
			data.WriteString(strings.TrimSpace(line[5:]))
		}
	}
}
