package middleware

import (
 "net/http"
 "net/http/httptest"
 "testing"
 "mail/database"
)

// Подготовим тестовые данные
func init() {
 database.UserHash = map[string]string{
  "valid_session": "user1",
 }
}

func TestAuthMiddleware(t *testing.T) {
 tests := []struct {
  name           string
  sessionCookie  string
  expectedStatus int
 }{
  {
   name:           "valid session",
   sessionCookie:  "valid_session",
   expectedStatus: http.StatusOK,
  },
  {
   name:           "invalid session",
   sessionCookie:  "invalid_session",
   expectedStatus: http.StatusUnauthorized,
  },
  {
   name:           "no session cookie",
   sessionCookie:  "",
   expectedStatus: http.StatusUnauthorized,
  },

 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   req := httptest.NewRequest(http.MethodGet, "/", nil)
   if tt.sessionCookie != "" {
    req.AddCookie(&http.Cookie{Name: "session", Value: tt.sessionCookie})
   } else if req.Method == http.MethodOptions {
    req.Method = http.MethodOptions
   }

   // Создаем ResponseRecorder для записи ответа
   rr := httptest.NewRecorder()

   // Оборачиваем наш обработчик в AuthMiddleware
   handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
   }))

   handler.ServeHTTP(rr, req)

   if status := rr.Code; status != tt.expectedStatus {
    t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
   }
  })
}
}