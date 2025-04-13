package middleware

import (
 "mail/config"
 "net/http"
 "net/http/httptest"
 "testing"
)

func TestCORS(t *testing.T) {
 // Настройка конфигурации для теста
 cfg := &config.Config{
  HTTPServer: struct {
		IP               string   `yaml:"ip"`
		Port             string   `yaml:"port"`
		AllowedIPsByCORS []string `yaml:"allowed_ips_by_cors"`
	}{
   AllowedIPsByCORS: []string{"http://localhost:4201"},
  },
 }

 // Создание тестового обработчика
 handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
 })

 // Обертывание тестового обработчика в middleware CORS
 corsHandler := CORS(handler, cfg)

 // Создание тестового запроса
 req, err := http.NewRequest(http.MethodOptions, "/", nil)
 if err != nil {
  t.Fatalf("Не удалось создать запрос: %v", err)
 }

 // Создание записывающего HTTP-ответа
 rr := httptest.NewRecorder()
 corsHandler.ServeHTTP(rr, req)

 // Проверка статуса ответа
 if status := rr.Code; status != http.StatusOK {
  t.Errorf("Неверный статус ответа: ожидали %v, получили %v", http.StatusOK, status)
 }

 // Проверка заголовков
 if origin := rr.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:4201" {
  t.Errorf("Неверный заголовок Access-Control-Allow-Origin: ожидали %v, получили %v", "http://localhost:4201", origin)
 }
 if methods := rr.Header().Get("Access-Control-Allow-Methods"); methods != "GET, POST, OPTIONS, DELETE" {
  t.Errorf("Неверный заголовок Access-Control-Allow-Methods: ожидали %v, получили %v", "GET, POST, OPTIONS, DELETE", methods)
 }
 if headers := rr.Header().Get("Access-Control-Allow-Headers"); headers != "Content-Type" {
  t.Errorf("Неверный заголовок Access-Control-Allow-Headers: ожидали %v, получили %v", "Content-Type", headers)
 }
 if credentials := rr.Header().Get("Access-Control-Allow-Credentials"); credentials != "true" {
  t.Errorf("Неверный заголовок Access-Control-Allow-Credentials: ожидали %v, получили %v", "true", credentials)
 }

 // Проверка обработки обычного запроса
 req, err = http.NewRequest(http.MethodGet, "/", nil)
 if err != nil {
  t.Fatalf("Не удалось создать запрос: %v", err)
 }
 rr = httptest.NewRecorder()
 corsHandler.ServeHTTP(rr, req)

 // Проверка статуса ответа для обычного запроса
 if status := rr.Code; status != http.StatusOK {
  t.Errorf("Неверный статус ответа: ожидали %v, получили %v", http.StatusOK, status)
 }
}