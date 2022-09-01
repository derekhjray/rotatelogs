package rotatelogs

import (
	"bytes"
	"encoding/json"
	isatty "github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	timeLayout = "2006/01/02 15:04:05.000"
)

func NewTextFormatter(disableColors bool) log.Formatter {
	return &textFormatter{
		indicators: map[log.Level]string{
			log.TraceLevel: "[T]",
			log.DebugLevel: "[D]",
			log.InfoLevel:  "[I]",
			log.WarnLevel:  "[W]",
			log.ErrorLevel: "[E]",
			log.FatalLevel: "[F]",
			log.PanicLevel: "[P]",
		},
		colors: map[log.Level]string{
			log.TraceLevel: "\033[32m",
			log.DebugLevel: "\033[32m",
			log.InfoLevel:  "\033[36m",
			log.WarnLevel:  "\033[33m",
			log.ErrorLevel: "\033[31m",
			log.FatalLevel: "\033[35m",
			log.PanicLevel: "\033[37m",
		},
		disableColors: disableColors,
	}
}

type textFormatter struct {
	indicators    map[log.Level]string
	colors        map[log.Level]string
	disableColors bool
}

func (tf *textFormatter) Format(entry *log.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = bytes.NewBuffer(make([]byte, 0, len(entry.Message)+32+len(entry.Data)*16))
	}

	buf.WriteString(entry.Time.Format(timeLayout))
	buf.WriteByte(' ')

	term := false
	if !tf.disableColors {
		if file, ok := entry.Logger.Out.(*os.File); ok && isatty.IsTerminal(file.Fd()) {
			term = true
		}
	}

	if term {
		buf.WriteString(tf.colors[entry.Level])
	}

	buf.WriteString(tf.indicators[entry.Level])

	if len(entry.Data) > 0 {
		if data, err := json.Marshal(entry.Data); err == nil {
			buf.WriteByte(' ')
			buf.Write(data)
		}
	}

	buf.WriteByte(' ')
	buf.WriteString(entry.Message)
	if term {
		buf.WriteString("\033[0m")
	}

	buf.WriteByte('\n')

	return buf.Bytes(), nil
}
