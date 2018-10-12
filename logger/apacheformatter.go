package logger

//格式如下[LEVEL] [2018-09-09|12:22:13] [module] [line] [msg]
import (
	"github.com/sirupsen/logrus"
	"fmt"
	"time"
	"bytes"
	"sync"
	"sort"
	"strings"
	"github.com/toney-li/go-common/util"
	"runtime"
)

const (
	nocolor                = 0
	red                    = 31
	green                  = 32
	yellow                 = 33
	blue                   = 36
	gray                   = 37
	defaultTimestampFormat = "2006-01-02|15:04:05"
	FieldKeyModule         = "module"
	FieldKeyLine           = "Line"
)

var (
	baseTimestamp time.Time
	emptyFieldMap logrus.FieldMap
)

func init() {
	baseTimestamp = time.Now()
}

type fieldKey string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[fieldKey]string

// TextFormatter formats logs into text
type ApacheFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &TextFormatter{
	//     FieldMap: FieldMap{
	//         FieldKeyTime:  "@timestamp",
	//         FieldKeyLevel: "@level",
	//         FieldKeyMsg:   "@message"}}
	FieldMap FieldMap

	sync.Once
}

func (f *ApacheFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = util.CheckIfTerminal(entry.Logger.Out)
	}
}

// Format renders a single log entry
func (f *ApacheFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	prefixFieldClashes(entry.Data, f.FieldMap)

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.Do(func() { f.init(entry) })

	isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if isColored {
		f.printColored(b, entry, keys, timestampFormat)
	} else {
		f.appendKeyValue(b, f.FieldMap.resolve(logrus.FieldKeyLevel), strings.ToUpper(entry.Level.String()))
		if !f.DisableTimestamp {
			f.appendKeyValue(b, f.FieldMap.resolve(logrus.FieldKeyTime), entry.Time.Format(timestampFormat))
		}
		pc, _, line, ok := runtime.Caller(6)
		if ok {
			f.appendKeyValue(b, f.FieldMap.resolve(FieldKeyModule), runtime.FuncForPC(pc).Name())
			f.appendKeyValue(b, FieldKeyLine, line)
		}
		if entry.Message != "" {
			f.appendKeyValue(b, f.FieldMap.resolve(logrus.FieldKeyMsg), entry.Message)
		}

		for _, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key])
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *ApacheFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation {
		levelText = levelText[0:4]
	}

	if f.DisableTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m %-44s ", levelColor, levelText, entry.Message)
	} else if !f.FullTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d] %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message)
	} else {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s] %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), entry.Message)
	}
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}
}

func (f *ApacheFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *ApacheFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString("[")
	//if key != logrus.FieldKeyLevel && key != logrus.FieldKeyTime {
	//	b.WriteString(key)
	//	b.WriteString("=")
	//}
	f.appendValue(b, value)
	b.WriteString("]")
}

func (f *ApacheFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	//if !f.needsQuoting(stringVal) {
	b.WriteString(stringVal)
	//} else {
	//	b.WriteString(fmt.Sprintf("%q", stringVal))
	//}
}

func prefixFieldClashes(data logrus.Fields, fieldMap FieldMap) {
	timeKey := fieldMap.resolve(logrus.FieldKeyTime)
	if t, ok := data[timeKey]; ok {
		data["fields."+timeKey] = t
	}

	msgKey := fieldMap.resolve(logrus.FieldKeyMsg)
	if m, ok := data[msgKey]; ok {
		data["fields."+msgKey] = m
	}

	levelKey := fieldMap.resolve(logrus.FieldKeyLevel)
	if l, ok := data[levelKey]; ok {
		data["fields."+levelKey] = l
	}
}

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}
	return string(key)
}
