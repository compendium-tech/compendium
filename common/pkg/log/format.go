package log

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var timeStampFormat = "2006-01-02T15:04:05.000000Z07:00"

type LogFormatter struct {
	Program     string
	Environment string
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+7)

	for k, v := range entry.Data {
		data[k] = v
	}

	data["time"] = entry.Time.UTC().Format(timeStampFormat)
	data["msg"] = entry.Message
	data["level"] = strings.ToUpper(entry.Level.String())
	data["program"] = f.Program
	data["environment"] = f.Environment

	if entry.HasCaller() {
		data["caller"] = entry.Caller.Function
		data["file"] = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
	}

	serialized, err := json.Marshal(data)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}

	return append(serialized, '\n'), nil
}
