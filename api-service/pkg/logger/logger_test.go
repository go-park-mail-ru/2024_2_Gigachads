package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"mail/api-service/internal/delivery/httpserver/email/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testBuffer struct {
	buf bytes.Buffer
}

func (b *testBuffer) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

func (b *testBuffer) String() string {
	return b.buf.String()
}

func TestLogger(t *testing.T) {
	buf := &testBuffer{}
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	testLogger := slog.New(slog.NewJSONHandler(buf, opts))
	logger := Logger{testLogger}

	tests := []struct {
		name    string
		logFunc func(msg string, args ...any)
		level   string
		message string
		args    []any
	}{
		{
			name:    "Info logging",
			logFunc: logger.Info,
			level:   "INFO",
			message: "test info message",
			args:    []any{"key", "value"},
		},
		{
			name:    "Debug logging",
			logFunc: logger.Debug,
			level:   "DEBUG",
			message: "test debug message",
			args:    []any{"debug_key", "debug_value"},
		},
		{
			name:    "Warn logging",
			logFunc: logger.Warn,
			level:   "WARN",
			message: "test warn message",
			args:    []any{"warn_key", "warn_value"},
		},
		{
			name:    "Error logging",
			logFunc: logger.Error,
			level:   "ERROR",
			message: "test error message",
			args:    []any{"error_key", "error_value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.buf.Reset()

			if len(tt.args) > 0 {
				tt.logFunc(tt.message, tt.args...)
			} else {
				tt.logFunc(tt.message)
			}

			var logEntry map[string]interface{}
			err := json.Unmarshal(buf.buf.Bytes(), &logEntry)
			assert.NoError(t, err)

			assert.Equal(t, tt.level, logEntry["level"])

			assert.Equal(t, tt.message, logEntry["msg"])

			if len(tt.args) > 0 {
				assert.Equal(t, tt.args[1], logEntry[tt.args[0].(string)])
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.Logger)

	assert.IsType(t, &slog.Logger{}, logger.Logger)
}

func TestLoggerMethods(t *testing.T) {
	buf := &testBuffer{}
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	testLogger := slog.New(slog.NewJSONHandler(buf, opts))
	logger := Logger{testLogger}

	testCases := []struct {
		name    string
		logFunc func()
	}{
		{"Info", func() { logger.Info("info message") }},
		{"Debug", func() { logger.Debug("debug message") }},
		{"Warn", func() { logger.Warn("warn message") }},
		{"Error", func() { logger.Error("error message") }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.buf.Reset()
			tc.logFunc()
			assert.NotEmpty(t, buf.String(), "Log output should not be empty")
		})
	}
}

func TestMockLogger_Debug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)

	mockLogger.EXPECT().Debug("Debug message").Times(1)

	mockLogger.Debug("Debug message")
}

func TestMockLogger_Warn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)

	mockLogger.EXPECT().Warn("Warning message").Times(1)

	mockLogger.Warn("Warning message")
}

func TestMockLogger_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogable(ctrl)

	mockLogger.EXPECT().Error("Error message").Times(1)

	mockLogger.Error("Error message")
}
