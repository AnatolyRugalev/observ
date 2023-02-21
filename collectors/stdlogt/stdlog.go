package stdlogt

import (
	"github.com/AnatolyRugalev/observ/logcollect"
	"github.com/AnatolyRugalev/observ/logq"
	"github.com/AnatolyRugalev/observ/logt"
	"io"
	"log"
	"strings"
	"sync/atomic"
	"time"
)

type Collector struct {
	sink   atomic.Pointer[logcollect.Sink]
	logger *log.Logger
	prefix string
	flags  int
	writer io.Writer
}

func (c *Collector) Write(p []byte) (n int, err error) {
	if sink := c.sink.Load(); sink != nil {
		// TODO: performance and safety
		record := logq.Record{
			Attributes: map[string]any{},
		}
		str := string(p)
		if c.flags&log.Lmsgprefix == 0 {
			str = str[len(c.prefix):]
		}
		needExtraSpace := false
		dateFormat := ""
		if c.flags&log.Ldate != 0 {
			dateFormat += "2006/01/02"
		}
		if c.flags&log.Ltime != 0 {
			dateFormat += " 15:04:05"
			needExtraSpace = true
		}
		if c.flags&log.Lmicroseconds != 0 {
			dateFormat += ".999999"
		}
		if dateFormat != "" {
			idx := strings.IndexFunc(str, func(r rune) bool {
				if r == ' ' {
					if needExtraSpace {
						needExtraSpace = false
						return false
					}
					return true
				}
				return false
			})
			if idx > 0 {
				t, err := time.Parse(dateFormat, str[:idx])
				if err != nil {
					panic(err)
				}
				record.Time = t
				str = str[idx+1:]
			}
		}
		if c.flags&log.Lshortfile != 0 || c.flags&log.Llongfile != 0 {
			idx := strings.IndexRune(str, ':')
			file := str[:idx]
			str = str[idx+1:]
			idx = strings.IndexRune(str, ':')
			line := str[:idx]
			str = str[idx+1:]
			record.Attributes[logq.AttributeFile] = file
			record.Attributes[logq.AttributeLine] = line
		}

		if c.flags&log.Lmsgprefix != 0 {
			str = str[len(c.prefix)+1:]
		}
		record.Message = str[:len(str)-1]
		(*sink).Add(record)
	}
	return c.writer.Write(p)
}

func (c *Collector) CaptureLogs(sink logcollect.Sink) (done func()) {
	old := c.sink.Swap(&sink)
	return func() {
		c.sink.Swap(old)
	}
}

func New(logger *log.Logger) logcollect.Collector {
	c := &Collector{
		logger: logger,
		prefix: logger.Prefix(),
		flags:  logger.Flags(),
		writer: logger.Writer(),
	}
	logger.SetOutput(c)
	return c
}

func (c *Collector) Close() {
	c.logger.SetOutput(c.writer)
}

// Logger sets up collector with specified log.Logger instance .
func Logger(logger *log.Logger) logt.Option {
	return logt.WithCollector(New(logger))
}

// Default returns collector for package-level logger (log.Default).
// If you want to capture all `log` package level logs, use this option.
func Default() logt.Option {
	return Logger(log.Default())
}
