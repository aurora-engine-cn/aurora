package web

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
)

// Log 自定义Log需要实现的借口
type Log interface {
	Info(...interface{})
	Error(...interface{})
	Debug(...interface{})
	Panic(...interface{})
	Warn(...interface{})
}
type Formatter struct {
	*logrus.TextFormatter
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	format := ""
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = &bytes.Buffer{}
	}
	t := entry.Time.Format("2006-01-02 15:04:05")
	if entry.Data != nil && len(entry.Data) > 0 {
		format = fmt.Sprintf("%s [%s] [%s]-> %s\n", t, entry.Level, entry.Data["type"], entry.Message)
	} else {
		format = fmt.Sprintf("%s [%s] -> %s\n", t, entry.Level, entry.Message)
	}
	buf.WriteString(format)
	return buf.Bytes(), nil
}
