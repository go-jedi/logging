package logging

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
)

const (
	defaultLevel      = LevelInfo
	defaultAddSource  = true
	defaultIsJSON     = true
	defaultSetFile    = false
	defaultSetDefault = true
)

func NewLogger(opts ...LoggerOption) *Logger {
	config := &LoggerOptions{
		Level:      defaultLevel,
		AddSource:  defaultAddSource,
		IsJSON:     defaultIsJSON,
		SetFile:    defaultSetFile,
		LogFile:    LogFile{},
		SetDefault: defaultSetDefault,
	}

	for _, opt := range opts {
		opt(config)
	}

	options := &HandlerOptions{
		AddSource: config.AddSource,
		Level:     config.Level,
	}

	var h Handler = NewTextHandler(os.Stdout, options)

	if config.IsJSON {
		h = NewJSONHandler(os.Stdout, options)
	}

	if config.SetFile {
		mw := io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   config.LogFile.Filename, // path to file
				MaxSize:    config.LogFile.MaxSize,  // megabytes
				MaxBackups: config.LogFile.MaxBackups,
				MaxAge:     config.LogFile.MaxAge, // days
			},
		)

		if config.IsJSON {
			h = NewJSONHandler(mw, options)
		} else {
			h = NewTextHandler(mw, options)
		}
	}

	logger := New(h)

	if config.SetDefault {
		SetDefault(logger)
	}

	return logger
}

type LogFile struct {
	Filename string
	// megabytes
	MaxSize int
	// max backups
	MaxBackups int
	// days
	MaxAge int
}

type LoggerOptions struct {
	Level      Level
	AddSource  bool
	IsJSON     bool
	SetFile    bool
	LogFile    LogFile
	SetDefault bool
}

type LoggerOption func(*LoggerOptions)

// WithLevel logger option sets the log level, if not set, the default level is Info.
// debug, info, warn, error.
func WithLevel(level string) LoggerOption {
	return func(o *LoggerOptions) {
		var l Level
		if err := l.UnmarshalText([]byte(level)); err != nil {
			l = LevelInfo
		}

		o.Level = l
	}
}

// WithAddSource logger option sets the add source option, which will add source file and line number to the log record.
func WithAddSource(addSource bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.AddSource = addSource
	}
}

// WithIsJSON logger option sets the is json option, which will set JSON format for the log record.
func WithIsJSON(isJSON bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.IsJSON = isJSON
	}
}

// WithSetFile logger option sets the set file option, which will set the created logger in file.
func WithSetFile(setFile bool, logFile LogFile) LoggerOption {
	return func(o *LoggerOptions) {
		o.SetFile = setFile
		o.LogFile = logFile
	}
}

// WithSetDefault logger option sets the set default option, which will set the created logger as default logger.
func WithSetDefault(setDefault bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.SetDefault = setDefault
	}
}

// WithAttrs returns logger with attributes.
func WithAttrs(ctx context.Context, attrs ...Attr) *Logger {
	logger := L(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

// WithDefaultAttrs returns logger with default attributes.
func WithDefaultAttrs(logger *Logger, attrs ...Attr) *Logger {
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

func L(ctx context.Context) *Logger {
	return loggerFromContext(ctx)
}

func Default() *Logger {
	return slog.Default()
}
