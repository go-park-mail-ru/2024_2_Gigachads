package logger

import(
	"log/slog"
	"os"
)

type Logable interface {
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type Logger struct{
	Logger Logable
}

func NewLogger() Logger{
	f, err := os.OpenFile("log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	logger := slog.New(slog.NewJSONHandler(f, nil)) //os.Stdout
	if err != nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.Error(err.Error())
	}
	
	//defer f.Close()
	//logger := slog.New(slog.NewJSONHandler(f, nil)) //os.Stdout
    return Logger{logger}
}

func (l Logger) Info(msg string, args ...any){
	l.Logger.Info(msg, args...)
}

func (l Logger) Debug(msg string, args ...any){
	l.Logger.Debug(msg, args...)
}

func (l Logger) Warn(msg string, args ...any){
	l.Logger.Warn(msg, args...)
}

func (l Logger) Error(msg string, args ...any){
	l.Logger.Error(msg, args...)
}