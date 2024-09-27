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
	{
		ID:       "4",
		From:     "michael.jones@example.com",
		To:       "jane.smith@example.com",
		Body:     "Good morning, Jane. I hope you're doing well. Can we schedule a quick call later today?",
		Title:    "Quick Call Request",
		Status:   "sent",
		Datetime: time.Now().Add(-72 * time.Hour),
	},
	{
		ID:       "5",
		From:     "anna.bell@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hey Jane, I’ve completed the task you assigned. Let me know if there’s anything else.",
		Title:    "Task Completion",
		Status:   "sent",
		Datetime: time.Now().Add(-6 * time.Hour),
	},
	{
		ID:       "6",
		From:     "peter.parker@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, just a heads up that I'll be on leave next week.",
		Title:    "Upcoming Leave Notice",
		Status:   "sent",
		Datetime: time.Now().Add(-48 * time.Hour),
	},
	{
		ID:       "7",
		From:     "david.morris@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, attached are the documents you requested. Let me know if everything is in order.",
		Title:    "Requested Documents",
		Status:   "sent",
		Datetime: time.Now().Add(-12 * time.Hour),
	},
	{
		ID:       "8",
		From:     "carla.jackson@example.com",
		To:       "jane.smith@example.com",
		Body:     "Jane, could you please confirm the time for tomorrow's meeting?",
		Title:    "Meeting Time Confirmation",
		Status:   "sent",
		Datetime: time.Now().Add(-30 * time.Hour),
	},
	{
		ID:       "9",
		From:     "william.turner@example.com",
		To:       "jane.smith@example.com",
		Body:     "Jane, please review the attached report before our next meeting.",
		Title:    "Report Review Request",
		Status:   "sent",
		Datetime: time.Now().Add(-20 * time.Hour),
	},
	{
		ID:       "10",
		From:     "emma.clark@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, can you send over the final budget for this project?",
		Title:    "Final Budget Request",
		Status:   "sent",
		Datetime: time.Now().Add(-10 * time.Hour),
	},
	{
		ID:       "11",
		From:     "chris.evans@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, please confirm if the design draft I sent is okay for you.",
		Title:    "Design Draft Confirmation",
		Status:   "sent",
		Datetime: time.Now().Add(-36 * time.Hour),
	},
	{
		ID:       "12",
		From:     "megan.jones@example.com",
		To:       "jane.smith@example.com",
		Body:     "Jane, could you review the contract and get back to me by Friday?",
		Title:    "Contract Review",
		Status:   "sent",
		Datetime: time.Now().Add(-12 * time.Hour),
	},
	{
		ID:       "13",
		From:     "joshua.kim@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hey Jane, don't forget to submit the expense report before the end of the week.",
		Title:    "Expense Report Reminder",
		Status:   "sent",
		Datetime: time.Now().Add(-5 * time.Hour),
	},
	{
		ID:       "14",
		From:     "emily.hall@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, I’ll be sending over the updated timeline later today. Let me know if any changes are needed.",
		Title:    "Updated Timeline Coming Soon",
		Status:   "sent",
		Datetime: time.Now().Add(-8 * time.Hour),
	},
	{
		ID:       "15",
		From:     "james.king@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, here is the proposal we discussed last week. Looking forward to your feedback.",
		Title:    "Proposal for Review",
		Status:   "sent",
		Datetime: time.Now().Add(-16 * time.Hour),
	},
	{
		ID:       "16",
		From:     "olivia.green@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, please let me know if you have received the shipment for the upcoming event.",
		Title:    "Shipment Confirmation",
		Status:   "sent",
		Datetime: time.Now().Add(-40 * time.Hour),
	},
	{
		ID:       "17",
		From:     "daniel.lee@example.com",
		To:       "jane.smith@example.com",
		Body:     "Jane, just a quick note to say that the presentation files are ready. I’ll send them over shortly.",
		Title:    "Presentation Files Ready",
		Status:   "sent",
		Datetime: time.Now().Add(-24 * time.Hour),
	},
	{
		ID:       "18",
		From:     "sophia.brown@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, can we reschedule our call for next Monday? Let me know if that works for you.",
		Title:    "Call Reschedule Request",
		Status:   "sent",
		Datetime: time.Now().Add(-15 * time.Hour),
	},
	{
		ID:       "19",
		From:     "liam.thomas@example.com",
		To:       "jane.smith@example.com",
		Body:     "Jane, just wanted to remind you about the client meeting at 3 PM tomorrow.",
		Title:    "Client Meeting Reminder",
		Status:   "sent",
		Datetime: time.Now().Add(-9 * time.Hour),
	},
	{
		ID:       "20",
		From:     "isabella.moore@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, attached is the final version of the marketing plan. Let me know if everything looks good.",
		Title:    "Final Marketing Plan",
		Status:   "sent",
		Datetime: time.Now().Add(-7 * time.Hour),
	},
	{
		ID:       "21",
		From:     "john.martin@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, just a quick reminder to submit your report by the end of the day.",
		Title:    "Report Submission Reminder",
		Status:   "sent",
		Datetime: time.Now().Add(-2 * time.Hour),
	},
	{
		ID:       "22",
		From:     "lucas.wilson@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hey Jane, I’ve updated the database as requested. Let me know if you need further changes.",
		Title:    "Database Update Completed",
		Status:   "sent",
		Datetime: time.Now().Add(-5 * time.Hour),
	},
	{
		ID:       "23",
		From:     "grace.evans@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, please find the attached travel itinerary for next week's conference.",
		Title:    "Conference Itinerary",
		Status:   "sent",
		Datetime: time.Now().Add(-18 * time.Hour),
	},
	{
		ID:       "24",
		From:     "ethan.james@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, just confirming that we received your feedback on the draft proposal. Thank you!",
		Title:    "Feedback Confirmation",
		Status:   "sent",
		Datetime: time.Now().Add(-12 * time.Hour),
	},
	{
		ID:       "25",
		From:     "ava.thomas@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hey Jane, we’ve rescheduled the team meeting to next Tuesday at 2 PM. Please confirm if that works for you.",
		Title:    "Team Meeting Reschedule",
		Status:   "sent",
		Datetime: time.Now().Add(-7 * time.Hour),
	},
	{
		ID:       "26",
		From:     "noah.robinson@example.com",
		To:       "jane.smith@example.com",
		Body:     "Jane, here is the finalized version of the project timeline. Let me know if everything looks good.",
		Title:    "Finalized Project Timeline",
		Status:   "sent",
		Datetime: time.Now().Add(-48 * time.Hour),
	},
	{
		ID:       "27",
		From:     "mia.harris@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, the client has confirmed the meeting time for Thursday at 11 AM. Please prepare the presentation.",
		Title:    "Client Meeting Confirmation",
		Status:   "sent",
		Datetime: time.Now().Add(-24 * time.Hour),
	},
	{
		ID:       "28",
		From:     "logan.lewis@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hey Jane, I’ve shared the latest budget report with you. Please review and send your feedback.",
		Title:    "Budget Report Review",
		Status:   "sent",
		Datetime: time.Now().Add(-6 * time.Hour),
	},
	{
		ID:       "29",
		From:     "ella.walker@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, the event planning checklist has been updated. Please take a look and let me know if anything is missing.",
		Title:    "Event Planning Checklist",
		Status:   "sent",
		Datetime: time.Now().Add(-4 * time.Hour),
	},
	{
		ID:       "30",
		From:     "ben.jackson@example.com",
		To:       "jane.smith@example.com",
		Body:     "Jane, the technical report you requested is ready. Let me know when you'd like to go over it.",
		Title:    "Technical Report Ready",
		Status:   "sent",
		Datetime: time.Now().Add(-20 * time.Hour),
	},
	{
		ID:       "31",
		From:     "sophia.jones@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, the marketing team has finalized the campaign strategy. Please review and provide feedback.",
		Title:    "Marketing Strategy Finalized",
		Status:   "sent",
		Datetime: time.Now().Add(-8 * time.Hour),
	},
	{
		ID:       "32",
		From:     "jack.smith@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, I’ve uploaded the final draft of the article. Please review it at your earliest convenience.",
		Title:    "Article Draft Review",
		Status:   "sent",
		Datetime: time.Now().Add(-14 * time.Hour),
	},
	{
		ID:       "33",
		From:     "chloe.davis@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, just a reminder to check the test results in the shared folder.",
		Title:    "Test Results Reminder",
		Status:   "sent",
		Datetime: time.Now().Add(-10 * time.Hour),
	},
	{
		ID:       "34",
		From:     "oliver.brown@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, I’ve updated the presentation with the latest data. Let me know if you need further adjustments.",
		Title:    "Presentation Updated",
		Status:   "sent",
		Datetime: time.Now().Add(-9 * time.Hour),
	},
	{
		ID:       "35",
		From:     "amelia.johnson@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, the venue for the event has been confirmed. Let’s discuss the next steps tomorrow.",
		Title:    "Event Venue Confirmation",
		Status:   "sent",
		Datetime: time.Now().Add(-18 * time.Hour),
	},
	{
		ID:       "36",
		From:     "liam.king@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, please review the attached project charter and let me know your thoughts.",
		Title:    "Project Charter Review",
		Status:   "sent",
		Datetime: time.Now().Add(-36 * time.Hour),
	},
	{
		ID:       "37",
		From:     "emma.scott@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, I’ve updated the design based on your feedback. Please check it and let me know if any further changes are needed.",
		Title:    "Design Update Based on Feedback",
		Status:   "sent",
		Datetime: time.Now().Add(-7 * time.Hour),
	},
	{
		ID:       "38",
		From:     "lucas.green@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, attached is the agenda for tomorrow’s meeting. Please review it beforehand.",
		Title:    "Meeting Agenda",
		Status:   "sent",
		Datetime: time.Now().Add(-12 * time.Hour),
	},
	{
		ID:       "39",
		From:     "charlotte.adams@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, I just wanted to confirm the deliverable deadlines for this quarter.",
		Title:    "Deliverable Deadlines Confirmation",
		Status:   "sent",
		Datetime: time.Now().Add(-24 * time.Hour),
	},
	{
		ID:       "40",
		From:     "henry.miller@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, the final draft of the policy document is ready for your approval.",
		Title:    "Policy Document Final Draft",
		Status:   "sent",
		Datetime: time.Now().Add(-3 * time.Hour),
	},
	{
		ID:       "41",
		From:     "daniel.clark@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, I’ve scheduled the call with the client for 2 PM on Friday. Let me know if that works for you.",
		Title:    "Client Call Scheduled",
		Status:   "sent",
		Datetime: time.Now().Add(-30 * time.Hour),
	},
	{
		ID:       "42",
		From:     "nora.taylor@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hey Jane, please find attached the updated version of the business plan. I’d love to hear your thoughts.",
		Title:    "Business Plan Update",
		Status:   "sent",
		Datetime: time.Now().Add(-20 * time.Hour),
	},
	{
		ID:       "43",
		From:     "jack.evans@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, could you provide an update on the budget approvals for the new project?",
		Title:    "Budget Approvals Update",
		Status:   "sent",
		Datetime: time.Now().Add(-10 * time.Hour),
	},
	{
		ID:       "44",
		From:     "maria.morris@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, the team is waiting for your input on the marketing materials. Please review them when you have a moment.",
		Title:    "Marketing Materials Review",
		Status:   "sent",
		Datetime: time.Now().Add(-8 * time.Hour),
	},
	{
		ID:       "45",
		From:     "david.young@example.com",
		To:       "jane.smith@example.com",
		Body:     "Jane, I’ve shared the preliminary data analysis with you. Please review and let me know if any adjustments are needed.",
		Title:    "Preliminary Data Analysis",
		Status:   "sent",
		Datetime: time.Now().Add(-4 * time.Hour),
	},
	{
		ID:       "46",
		From:     "elizabeth.carter@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, I’ve reserved the conference room for our meeting on Wednesday at 3 PM. Let me know if that works.",
		Title:    "Conference Room Reservation",
		Status:   "sent",
		Datetime: time.Now().Add(-12 * time.Hour),
	},
	{
		ID:       "47",
		From:     "william.lee@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, I’ve uploaded the video tutorial to the shared folder. Please review and confirm it’s all good.",
		Title:    "Video Tutorial Uploaded",
		Status:   "sent",
		Datetime: time.Now().Add(-22 * time.Hour),
	},
	{
		ID:       "48",
		From:     "emily.davis@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, just following up on the legal review. Have you had a chance to look at the document?",
		Title:    "Legal Review Follow-Up",
		Status:   "sent",
		Datetime: time.Now().Add(-5 * time.Hour),
	},
	{
		ID:       "49",
		From:     "michael.martinez@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, please find the attached risk assessment report. Let me know if you have any questions.",
		Title:    "Risk Assessment Report",
		Status:   "sent",
		Datetime: time.Now().Add(-3 * time.Hour),
	},
	{
		ID:       "50",
		From:     "sophia.thomas@example.com",
		To:       "jane.smith@example.com",
		Body:     "Hi Jane, just a reminder to review the client’s feedback on the recent design proposal.",
		Title:    "Client Feedback Reminder",
		Status:   "sent",
		Datetime: time.Now().Add(-6 * time.Hour),
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
