package logger

import(
	"log/slog"
	"os"
)

// type Logger struct{

// }

func NewLogger(){
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)
}

// func Info(str string){
// 	slog.Info(str)
// }

func Info(msg string, args ...any){
	slog.Info(msg, args...)
}

func Debug(msg string, args ...any){
	slog.Debug(msg, args...)
}

func Warn(msg string, args ...any){
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any){
	slog.Error(msg, args...)
}