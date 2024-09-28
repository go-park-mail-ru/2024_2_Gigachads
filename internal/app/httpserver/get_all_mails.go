package httpserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

type mockedTableStructure struct {
	Author      string    `json:"author"`      // UUID // UUID
	Description string    `json:"description"` // json
	Text        string    `json:"text"`
	Badge_text  string    `json:"badge_text"`
	Badge_type  string    `json:"badge_type"`
	Date        time.Time `json:"date"`
}

type Mails []mockedTableStructure

func (mails Mails) compare(otherMails Mails) bool {
	if len(mails) != len(otherMails) {
		return false
	}

	for i := range mails {
		if reflect.DeepEqual(otherMails[i], mails[i]) && compareDatetimes(otherMails[i].Date, mails[i].Date) {
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

		Author:      "john.doe@example.com",
		Description: "Hi Jane, just wanted to check in on the project status. Let me know if you need any help.",
		Text:        "Project Status Check-In",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-48 * time.Hour),
	},
	{

		Author:      "mark.brown@example.com",
		Description: "Hey Jane, just a reminder about our meeting tomorrow at 10 AM. Please confirm if you're available.",
		Text:        "Meeting Reminder",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-24 * time.Hour),
	},
	{

		Author:      "lisa.white@example.com",
		Description: "Hi Jane, I’ve attached the latest version of the report. Please review and send your feedback.",
		Text:        "Report Update",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now(),
	},
	{

		Author:      "michael.jones@example.com",
		Description: "Good morning, Jane. I hope you're doing well. Can we schedule a quick call later today?",
		Text:        "Quick Call Request",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-72 * time.Hour),
	},
	{

		Author:      "anna.bell@example.com",
		Description: "Hey Jane, I’ve completed the task you assigned. Let me know if there’s anything else.",
		Text:        "Task Completion",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-6 * time.Hour),
	},
	{

		Author:      "peter.parker@example.com",
		Description: "Hi Jane, just a heads up that I'll be on leave next week.",
		Text:        "Upcoming Leave Notice",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-48 * time.Hour),
	},
	{

		Author:      "david.morris@example.com",
		Description: "Hi Jane, attached are the documents you requested. Let me know if everything is in order.",
		Text:        "Requested Documents",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-12 * time.Hour),
	},
	{

		Author:      "carla.jackson@example.com",
		Description: "Jane, could you please confirm the time for tomorrow's meeting?",
		Text:        "Meeting Time Confirmation",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-30 * time.Hour),
	},
	{

		Author:      "william.turner@example.com",
		Description: "Jane, please review the attached report before our next meeting.",
		Text:        "Report Review Request",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-20 * time.Hour),
	},
	{

		Author:      "emma.clark@example.com",
		Description: "Hi Jane, can you send over the final budget for this project?",
		Text:        "Final Budget Request",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-10 * time.Hour),
	},
	{

		Author:      "chris.evans@example.com",
		Description: "Hi Jane, please confirm if the design draft I sent is okay for you.",
		Text:        "Design Draft Confirmation",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-36 * time.Hour),
	},
	{

		Author:      "megan.jones@example.com",
		Description: "Jane, could you review the contract and get back to me by Friday?",
		Text:        "Contract Review",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-12 * time.Hour),
	},
	{
		Author:      "joshua.kim@example.com",
		Description: "Hey Jane, don't forget to submit the expense report before the end of the week.",
		Text:        "Expense Report Reminder",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-5 * time.Hour),
	},
	{

		Author:      "emily.hall@example.com",
		Description: "Hi Jane, I’ll be sending over the updated timeline later today. Let me know if any changes are needed.",
		Text:        "Updated Timeline Coming Soon",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-8 * time.Hour),
	},
	{

		Author:      "james.king@example.com",
		Description: "Hi Jane, here is the proposal we discussed last week. Looking forward to your feedback.",
		Text:        "Proposal for Review",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-16 * time.Hour),
	},
	{

		Author:      "olivia.green@example.com",
		Description: "Hi Jane, please let me know if you have received the shipment for the upcoming event.",
		Text:        "Shipment Confirmation",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-40 * time.Hour),
	},
	{

		Author:      "daniel.lee@example.com",
		Description: "Jane, just a quick note to say that the presentation files are ready. I’ll send them over shortly.",
		Text:        "Presentation Files Ready",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-24 * time.Hour),
	},
	{

		Author:      "sophia.brown@example.com",
		Description: "Hi Jane, can we reschedule our call for next Monday? Let me know if that works for you.",
		Text:        "Call Reschedule Request",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-15 * time.Hour),
	},
	{

		Author:      "liam.thomas@example.com",
		Description: "Jane, just wanted to remind you about the client meeting at 3 PM tomorrow.",
		Text:        "Client Meeting Reminder",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-9 * time.Hour),
	},
	{

		Author:      "isabella.moore@example.com",
		Description: "Hi Jane, attached is the final version of the marketing plan. Let me know if everything looks good.",
		Text:        "Final Marketing Plan",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-7 * time.Hour),
	},
	{

		Author:      "john.martin@example.com",
		Description: "Hi Jane, just a quick reminder to submit your report by the end of the day.",
		Text:        "Report Submission Reminder",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-2 * time.Hour),
	},
	{

		Author:      "lucas.wilson@example.com",
		Description: "Hey Jane, I’ve updated the database as requested. Let me know if you need further changes.",
		Text:        "Database Update Completed",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-5 * time.Hour),
	},
	{

		Author:      "grace.evans@example.com",
		Description: "Hi Jane, please find the attached travel itinerary for next week's conference.",
		Text:        "Conference Itinerary",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-18 * time.Hour),
	},
	{

		Author:      "ethan.james@example.com",
		Description: "Hi Jane, just confirming that we received your feedback on the draft proposal. Thank you!",
		Text:        "Feedback Confirmation",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-12 * time.Hour),
	},
	{

		Author:      "ava.thomas@example.com",
		Description: "Hey Jane, we’ve rescheduled the team meeting to next Tuesday at 2 PM. Please confirm if that works for you.",
		Text:        "Team Meeting Reschedule",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-7 * time.Hour),
	},
	{

		Author:      "noah.robinson@example.com",
		Description: "Jane, here is the finalized version of the project timeline. Let me know if everything looks good.",
		Text:        "Finalized Project Timeline",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-48 * time.Hour),
	},
	{

		Author:      "mia.harris@example.com",
		Description: "Hi Jane, the client has confirmed the meeting time for Thursday at 11 AM. Please prepare the presentation.",
		Text:        "Client Meeting Confirmation",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-24 * time.Hour),
	},
	{

		Author:      "logan.lewis@example.com",
		Description: "Hey Jane, I’ve shared the latest budget report with you. Please review and send your feedback.",
		Text:        "Budget Report Review",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-6 * time.Hour),
	},
	{

		Author:      "ella.walker@example.com",
		Description: "Hi Jane, the event planning checklist has been updated. Please take a look and let me know if anything is missing.",
		Text:        "Event Planning Checklist",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-4 * time.Hour),
	},
	{
		Author:      "ben.jackson@example.com",
		Description: "Jane, the technical report you requested is ready. Let me know when you'd like to go over it.",
		Text:        "Technical Report Ready",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-20 * time.Hour),
	},
	{
		Author:      "sophia.jones@example.com",
		Description: "Hi Jane, the marketing team has finalized the campaign strategy. Please review and provide feedback.",
		Text:        "Marketing Strategy Finalized",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-8 * time.Hour),
	},
	{
		Author:      "jack.smith@example.com",
		Description: "Hi Jane, I’ve uploaded the final draft of the article. Please review it at your earliest convenience.",
		Text:        "Article Draft Review",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-14 * time.Hour),
	},
	{

		Author:      "chloe.davis@example.com",
		Description: "Hi Jane, just a reminder to check the test results in the shared folder.",
		Text:        "Test Results Reminder",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-10 * time.Hour),
	},
	{

		Author:      "oliver.brown@example.com",
		Description: "Hi Jane, I’ve updated the presentation with the latest data. Let me know if you need further adjustments.",
		Text:        "Presentation Updated",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-9 * time.Hour),
	},
	{

		Author:      "amelia.johnson@example.com",
		Description: "Hi Jane, the venue for the event has been confirmed. Let’s discuss the next steps tomorrow.",
		Text:        "Event Venue Confirmation",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-18 * time.Hour),
	},
	{

		Author:      "liam.king@example.com",
		Description: "Hi Jane, please review the attached project charter and let me know your thoughts.",
		Text:        "Project Charter Review",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-36 * time.Hour),
	},
	{

		Author:      "emma.scott@example.com",
		Description: "Hi Jane, I’ve updated the design based on your feedback. Please check it and let me know if any further changes are needed.",
		Text:        "Design Update Based on Feedback",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-7 * time.Hour),
	},
	{

		Author:      "lucas.green@example.com",
		Description: "Hi Jane, attached is the agenda for tomorrow’s meeting. Please review it beforehand.",
		Text:        "Meeting Agenda",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-12 * time.Hour),
	},
	{

		Author:      "charlotte.adams@example.com",
		Description: "Hi Jane, I just wanted to confirm the deliverable deadlines for this quarter.",
		Text:        "Deliverable Deadlines Confirmation",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-24 * time.Hour),
	},
	{

		Author:      "henry.miller@example.com",
		Description: "Hi Jane, the final draft of the policy document is ready for your approval.",
		Text:        "Policy Document Final Draft",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-3 * time.Hour),
	},
	{

		Author:      "daniel.clark@example.com",
		Description: "Hi Jane, I’ve scheduled the call with the client for 2 PM on Friday. Let me know if that works for you.",
		Text:        "Client Call Scheduled",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-30 * time.Hour),
	},
	{

		Author:      "nora.taylor@example.com",
		Description: "Hey Jane, please find attached the updated version of the business plan. I’d love to hear your thoughts.",
		Text:        "Business Plan Update",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-20 * time.Hour),
	},
	{

		Author:      "jack.evans@example.com",
		Description: "Hi Jane, could you provide an update on the budget approvals for the new project?",
		Text:        "Budget Approvals Update",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-10 * time.Hour),
	},
	{

		Author:      "maria.morris@example.com",
		Description: "Hi Jane, the team is waiting for your input on the marketing materials. Please review them when you have a moment.",
		Text:        "Marketing Materials Review",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-8 * time.Hour),
	},
	{

		Author:      "david.young@example.com",
		Description: "Jane, I’ve shared the preliminary data analysis with you. Please review and let me know if any adjustments are needed.",
		Text:        "Preliminary Data Analysis",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-4 * time.Hour),
	},
	{

		Author:      "elizabeth.carter@example.com",
		Description: "Hi Jane, I’ve reserved the conference room for our meeting on Wednesday at 3 PM. Let me know if that works.",
		Text:        "Conference Room Reservation",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-12 * time.Hour),
	},
	{

		Author:      "william.lee@example.com",
		Description: "Hi Jane, I’ve uploaded the video tutorial to the shared folder. Please review and confirm it’s all good.",
		Text:        "Video Tutorial Uploaded",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-22 * time.Hour),
	},
	{

		Author:      "emily.davis@example.com",
		Description: "Hi Jane, just following up on the legal review. Have you had a chance to look at the document?",
		Text:        "Legal Review Follow-Up",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-5 * time.Hour),
	},
	{

		Author:      "michael.martinez@example.com",
		Description: "Hi Jane, please find the attached risk assessment report. Let me know if you have any questions.",
		Text:        "Risk Assessment Report",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-3 * time.Hour),
	},
	{

		Author:      "sophia.thomas@example.com",
		Description: "Hi Jane, just a reminder to review the client’s feedback on the recent design proposal.",
		Text:        "Client Feedback Reminder",
		Badge_text:  "sent",
		Badge_type:  "test",
		Date:        time.Now().Add(-6 * time.Hour),
	},
}

func getAllMails(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("user_id")
	if err != nil {
		slog.Error("cannot get user_id from cookie")
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
	userID := cookie.Value
	fmt.Println(userID)

	result := make(Mails, 0)
	result = append(result, mockedMails...)

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
}
