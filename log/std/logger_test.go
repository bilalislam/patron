package std

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/beatlabs/patron/log"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	var b bytes.Buffer
	logger := New(&b, log.InfoLevel, map[string]interface{}{"name": "john doe", "age": 18})
	assert.NotNil(t, logger.debug)
	assert.NotNil(t, logger.info)
	assert.NotNil(t, logger.warn)
	assert.NotNil(t, logger.error)
	assert.NotNil(t, logger.fatal)
	assert.NotNil(t, logger.panic)
	assert.Equal(t, log.InfoLevel, logger.Level())
	assert.Equal(t, logger.fields, map[string]interface{}{"name": "john doe", "age": 18})
	assert.Contains(t, logger.fieldsLine, "age=18")
	assert.Contains(t, logger.fieldsLine, "name=john doe")
}

func TestNewSub(t *testing.T) {
	var b bytes.Buffer
	logger := New(&b, log.InfoLevel, map[string]interface{}{"name": "john doe"})
	assert.NotNil(t, logger)
	subLogger := logger.Sub(map[string]interface{}{"age": 18}).(*Logger)
	assert.NotNil(t, subLogger.debug)
	assert.NotNil(t, subLogger.info)
	assert.NotNil(t, subLogger.warn)
	assert.NotNil(t, subLogger.error)
	assert.NotNil(t, subLogger.fatal)
	assert.NotNil(t, subLogger.panic)
	assert.Equal(t, log.InfoLevel, subLogger.Level())
	assert.Equal(t, subLogger.fields, map[string]interface{}{"name": "john doe", "age": 18})
	assert.Contains(t, subLogger.fieldsLine, "age=18")
	assert.Contains(t, subLogger.fieldsLine, "name=john doe")
}

func TestLogger(t *testing.T) {
	// BEWARE: Since we are testing the log output change in line number of statements affect the test outcome
	var b bytes.Buffer
	logger := New(&b, log.DebugLevel, map[string]interface{}{"name": "john doe", "age": 18})

	type args struct {
		lvl  log.Level
		msg  string
		args []interface{}
	}
	tests := map[string]struct {
		args args
	}{
		"debug":  {args: args{lvl: log.DebugLevel, args: []interface{}{"hello world"}}},
		"debugf": {args: args{lvl: log.DebugLevel, msg: "Hi, %s", args: []interface{}{"John"}}},
		"info":   {args: args{lvl: log.InfoLevel, args: []interface{}{"hello world"}}},
		"infof":  {args: args{lvl: log.InfoLevel, msg: "Hi, %s", args: []interface{}{"John"}}},
		"warn":   {args: args{lvl: log.WarnLevel, args: []interface{}{"hello world"}}},
		"warnf":  {args: args{lvl: log.WarnLevel, msg: "Hi, %s", args: []interface{}{"John"}}},
		"error":  {args: args{lvl: log.ErrorLevel, args: []interface{}{"hello world"}}},
		"errorf": {args: args{lvl: log.ErrorLevel, msg: "Hi, %s", args: []interface{}{"John"}}},
		"panic":  {args: args{lvl: log.PanicLevel, args: []interface{}{"hello world"}}},
		"panicf": {args: args{lvl: log.PanicLevel, msg: "Hi, %s", args: []interface{}{"John"}}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			defer b.Reset()

			switch tt.args.lvl {
			case log.DebugLevel:
				if tt.args.msg == "" {
					logger.Debug(tt.args.args...)
				} else {
					logger.Debugf(tt.args.msg, tt.args.args...)
				}
			case log.InfoLevel:
				if tt.args.msg == "" {
					logger.Info(tt.args.args...)
				} else {
					logger.Infof(tt.args.msg, tt.args.args...)
				}
			case log.WarnLevel:
				if tt.args.msg == "" {
					logger.Warn(tt.args.args...)
				} else {
					logger.Warnf(tt.args.msg, tt.args.args...)
				}
			case log.ErrorLevel:
				if tt.args.msg == "" {
					logger.Error(tt.args.args...)
				} else {
					logger.Errorf(tt.args.msg, tt.args.args...)
				}
			case log.PanicLevel:
				if tt.args.msg == "" {
					assert.Panics(t, func() {
						logger.Panic(tt.args.args...)
					})
				} else {
					assert.Panics(t, func() {
						logger.Panicf(tt.args.msg, tt.args.args...)
					})
				}
			}

			if tt.args.msg == "" {
				assert.Contains(t, b.String(), fmt.Sprintf("%s age=18 name=john doe hello world", levelMap[tt.args.lvl]))
			} else {
				assert.Contains(t, b.String(), fmt.Sprintf("%s age=18 name=john doe Hi, John", levelMap[tt.args.lvl]))
			}
		})
	}
}

func TestLogger_shouldLog(t *testing.T) {
	type args struct {
		lvl log.Level
	}
	tests := map[string]struct {
		setupLevel log.Level
		args       args
		want       bool
	}{
		"setup debug,passing debug":    {setupLevel: log.DebugLevel, args: args{lvl: log.DebugLevel}, want: true},
		"setup debug,passing info":     {setupLevel: log.DebugLevel, args: args{lvl: log.InfoLevel}, want: true},
		"setup debug,passing warn":     {setupLevel: log.DebugLevel, args: args{lvl: log.WarnLevel}, want: true},
		"setup debug,passing error":    {setupLevel: log.DebugLevel, args: args{lvl: log.ErrorLevel}, want: true},
		"setup debug,passing panic":    {setupLevel: log.DebugLevel, args: args{lvl: log.PanicLevel}, want: true},
		"setup debug,passing fatal":    {setupLevel: log.DebugLevel, args: args{lvl: log.FatalLevel}, want: true},
		"setup info,passing debug":     {setupLevel: log.InfoLevel, args: args{lvl: log.DebugLevel}, want: false},
		"setup info,passing info":      {setupLevel: log.InfoLevel, args: args{lvl: log.InfoLevel}, want: true},
		"setup info,passing warn":      {setupLevel: log.InfoLevel, args: args{lvl: log.WarnLevel}, want: true},
		"setup info,passing error":     {setupLevel: log.InfoLevel, args: args{lvl: log.ErrorLevel}, want: true},
		"setup info,passing panic":     {setupLevel: log.InfoLevel, args: args{lvl: log.PanicLevel}, want: true},
		"setup info,passing fatal":     {setupLevel: log.InfoLevel, args: args{lvl: log.FatalLevel}, want: true},
		"setup warn,passing debug":     {setupLevel: log.WarnLevel, args: args{lvl: log.DebugLevel}, want: false},
		"setup warn,passing info":      {setupLevel: log.WarnLevel, args: args{lvl: log.InfoLevel}, want: false},
		"setup warn,passing warn":      {setupLevel: log.WarnLevel, args: args{lvl: log.WarnLevel}, want: true},
		"setup warn,passing error":     {setupLevel: log.WarnLevel, args: args{lvl: log.ErrorLevel}, want: true},
		"setup warn,passing panic":     {setupLevel: log.WarnLevel, args: args{lvl: log.PanicLevel}, want: true},
		"setup warn,passing fatal":     {setupLevel: log.WarnLevel, args: args{lvl: log.FatalLevel}, want: true},
		"setup error,passing debug":    {setupLevel: log.ErrorLevel, args: args{lvl: log.DebugLevel}, want: false},
		"setup error,passing info":     {setupLevel: log.ErrorLevel, args: args{lvl: log.InfoLevel}, want: false},
		"setup error,passing warn":     {setupLevel: log.ErrorLevel, args: args{lvl: log.WarnLevel}, want: false},
		"setup error,passing error":    {setupLevel: log.ErrorLevel, args: args{lvl: log.ErrorLevel}, want: true},
		"setup error,passing panic":    {setupLevel: log.ErrorLevel, args: args{lvl: log.PanicLevel}, want: true},
		"setup error,passing fatal":    {setupLevel: log.ErrorLevel, args: args{lvl: log.FatalLevel}, want: true},
		"setup fatal,passing debug":    {setupLevel: log.FatalLevel, args: args{lvl: log.DebugLevel}, want: false},
		"setup fatal,passing info":     {setupLevel: log.FatalLevel, args: args{lvl: log.InfoLevel}, want: false},
		"setup fatal,passing warn":     {setupLevel: log.FatalLevel, args: args{lvl: log.WarnLevel}, want: false},
		"setup fatal,passing error":    {setupLevel: log.FatalLevel, args: args{lvl: log.ErrorLevel}, want: false},
		"setup fatal,passing fatal":    {setupLevel: log.FatalLevel, args: args{lvl: log.FatalLevel}, want: true},
		"setup fatal,passing panic":    {setupLevel: log.FatalLevel, args: args{lvl: log.PanicLevel}, want: true},
		"setup panic,passing debug":    {setupLevel: log.PanicLevel, args: args{lvl: log.DebugLevel}, want: false},
		"setup panic,passing info":     {setupLevel: log.PanicLevel, args: args{lvl: log.InfoLevel}, want: false},
		"setup panic,passing warn":     {setupLevel: log.PanicLevel, args: args{lvl: log.WarnLevel}, want: false},
		"setup panic,passing error":    {setupLevel: log.PanicLevel, args: args{lvl: log.ErrorLevel}, want: false},
		"setup panic,passing fatal":    {setupLevel: log.PanicLevel, args: args{lvl: log.FatalLevel}, want: false},
		"setup panic,passing panic":    {setupLevel: log.PanicLevel, args: args{lvl: log.PanicLevel}, want: true},
		"setup no level,passing debug": {setupLevel: log.NoLevel, args: args{lvl: log.DebugLevel}, want: false},
		"setup no level,passing info":  {setupLevel: log.NoLevel, args: args{lvl: log.InfoLevel}, want: false},
		"setup no level,passing warn":  {setupLevel: log.NoLevel, args: args{lvl: log.WarnLevel}, want: false},
		"setup no level,passing error": {setupLevel: log.NoLevel, args: args{lvl: log.ErrorLevel}, want: false},
		"setup no level,passing panic": {setupLevel: log.NoLevel, args: args{lvl: log.PanicLevel}, want: false},
		"setup no level,passing fatal": {setupLevel: log.NoLevel, args: args{lvl: log.FatalLevel}, want: false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			l := &Logger{level: tt.setupLevel}
			assert.Equal(t, tt.want, l.shouldLog(tt.args.lvl))
		})
	}
}

var buf bytes.Buffer

func BenchmarkLogger(b *testing.B) {
	var tmpBuf bytes.Buffer
	logger := New(&tmpBuf, log.DebugLevel, map[string]interface{}{"name": "john doe", "age": 18})
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Debugf("Hello %s!", "John")
	}
	buf = tmpBuf
}
