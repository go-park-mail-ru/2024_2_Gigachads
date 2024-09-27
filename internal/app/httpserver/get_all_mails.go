package httpserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

type errorResponse struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

type mockedTableStructure struct {
	ID       string    `json:"id"`   // UUID
	From     string    `json:"from"` // UUID
	To       string    `json:"to"`   // UUID
	Body     string    `json:"body"` // json
	Title    string    `json:"title"`
	Status   string    `json:"status"`
	Datetime time.Time `json:"datetime"`
}

type Mails []mockedTableStructure

func (mails Mails) compare(otherMails Mails) bool {
	if len(mails) != len(otherMails) {
		return false
	}

	for i := range mails {
		if reflect.DeepEqual(otherMails[i], mails[i]) && compareDatetimes(otherMails[i].Datetime, mails[i].Datetime) {
			return false
		}
	}
	return true
}

func compareDatetimes(timeOne time.Time, timeTwo time.Time) bool {
	return timeOne.Format("2024-09-23 22:51:51.816489625 +0300") == timeTwo.Format("2024-09-23 22:51:51.816489625 +0300")
}

var mockedMails = Mails{
	{
		ID:       "1",
		From:     "john.doe@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, just wanted to check in on the project status. Let me know if you need any help.",
		Title:    "Project Status Check-In",
		Status:   "sent",
		Datetime: time.Now().Add(-48 * time.Hour),
	},
	{
		ID:       "2",
		From:     "mark.brown@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hey Jane, just a reminder about our meeting tomorrow at 10 AM. Please confirm if you're available.",
		Title:    "Meeting Reminder",
		Status:   "sent",
		Datetime: time.Now().Add(-24 * time.Hour),
	},
	{
		ID:       "3",
		From:     "lisa.white@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, I’ve attached the latest version of the report. Please review and send your feedback.",
		Title:    "Report Update",
		Status:   "sent",
		Datetime: time.Now(),
	},
}

func getAllMails(w http.ResponseWriter, req *http.Request) {
	userID, ok := req.Context().Value("user-id").(string)
	fmt.Println(userID)
	if !ok {
		slog.Error("cannot type assertion userID to string")
		w.WriteHeader(http.StatusForbidden)
		response := errorResponse{
			Status: http.StatusForbidden,
			Body:   "Validation_error",
		}
		marshaledResponse, err := json.Marshal(response)
		if err != nil {
			slog.Error("failed to marshal error response")
		}
		w.Write(marshaledResponse)
		return
	}

	result := make(Mails, 0)
	for _, message := range mockedMails {
		if message.To == userID {
			result = append(result, message)
		}
	}

	resultToJson, err := json.Marshal(result)
	if err != nil {
		slog.Error(fmt.Sprintf("cannot convert to json: %v", err))
		response := errorResponse{
			Status: http.StatusInternalServerError,
			Body:   "Internal_error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		marshaledResponse, err := json.Marshal(response)
		if err != nil {
			slog.Error("failed to marshal error response")
		}
		w.Write(marshaledResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultToJson)
	return
}
