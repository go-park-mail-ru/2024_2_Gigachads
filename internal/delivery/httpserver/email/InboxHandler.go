package email

import (
	"encoding/json"
	"mail/pkg/utils"
	"net/http"
)

func (er *EmailRouter) InboxHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	mails, err := er.EmailUseCase.Inbox(cookie.Value)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := json.Marshal(mails)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "json error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
