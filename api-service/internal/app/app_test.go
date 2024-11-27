package app

import (
	"testing"
	"time"

	"mail/api-service/pkg/logger"
	"mail/config"
)

func TestRun(t *testing.T) {
	// Создаём тестовый конфиг и логгер
	cfg := &config.Config{}
	log := logger.NewLogger()

	stop := make(chan struct{})

	go func() {
		Run(cfg, log)
		close(stop)
	}()

	time.Sleep(100 * time.Millisecond)

	select {
	case <-stop:
		t.Error("приложение завершилось преждевременно")
	default:
	}
}
