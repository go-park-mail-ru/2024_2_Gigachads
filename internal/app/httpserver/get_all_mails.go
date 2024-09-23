package httpserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// TODO: после появления ручки создания письма в базе убрать
type mockedTableStructure struct {
	ID       string    `json:"id"`   // UUID
	From     string    `json:"from"` // UUID
	To       string    `json:"to"`   // UUID
	Body     string    `json:"body"` // json
	Title    string    `json:"title"`
	Status   string    `json:"status"`
	Datetime time.Time `json:"datetime"`
}

// TODO: после появления ручки создания письма в базе убрать
var mockedDatabase = make([]mockedTableStructure, 0)

// TODO: после появления ручки создания письма в базе убрать
var testMockedItem = mockedTableStructure{To: "test-uuid", Title: "test"}

// TODO: после появления ручки создания письма в базе убрать
func prepareMockedDatabase() {
	mockedDatabase = append(mockedDatabase, testMockedItem)
}

func getAllMails(w http.ResponseWriter, req *http.Request) {
	// TODO: после появления ручки создания письма в базе убрать
	prepareMockedDatabase()

	userID, ok := req.Context().Value("user-id").(string)
	if !ok {
		slog.Error("cannot type assertion userID to string")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	result := make([]mockedTableStructure, 0)
	for _, message := range mockedDatabase {
		if message.To == userID {
			result = append(result, message)
		}
	}

	resultToJson, err := json.Marshal(result)
	if err != nil {
		slog.Error(fmt.Sprintf("cannot convert to json: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resultToJson)
	return
}
