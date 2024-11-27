package smtp

import (
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSendEmailRequest(t *testing.T) {
	req := &SendEmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	t.Run("GetFrom", func(t *testing.T) {
		if got := req.GetFrom(); got != "sender@example.com" {
			t.Errorf("GetFrom() = %v, want %v", got, "sender@example.com")
		}
	})

	t.Run("GetTo", func(t *testing.T) {
		if got := req.GetTo(); got != "recipient@example.com" {
			t.Errorf("GetTo() = %v, want %v", got, "recipient@example.com")
		}
	})

	t.Run("GetSubject", func(t *testing.T) {
		if got := req.GetSubject(); got != "Test Subject" {
			t.Errorf("GetSubject() = %v, want %v", got, "Test Subject")
		}
	})

	t.Run("GetBody", func(t *testing.T) {
		if got := req.GetBody(); got != "Test Body" {
			t.Errorf("GetBody() = %v, want %v", got, "Test Body")
		}
	})
}

func TestReplyEmailRequestMethods(t *testing.T) {
	fixedTime := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	timestamp := timestamppb.New(fixedTime)

	req := &ReplyEmailRequest{
		From:        "sender@example.com",
		To:          "recipient@example.com",
		Title:       "Re: Test",
		Sender:      "original@example.com",
		ReplyText:   "Reply content",
		SendingDate: timestamp,
		Description: "Reply Description",
	}

	t.Run("Reset", func(t *testing.T) {
		reqCopy := *req
		reqCopy.Reset()

		if reqCopy.From != "" || reqCopy.To != "" || reqCopy.Title != "" ||
			reqCopy.Sender != "" || reqCopy.ReplyText != "" ||
			reqCopy.Description != "" || reqCopy.SendingDate != nil {
			t.Error("Reset() did not clear all fields")
		}
	})

	t.Run("String", func(t *testing.T) {
		str := req.String()
		expectedFields := []string{
			"sender@example.com",
			"recipient@example.com",
			"Re: Test",
			"original@example.com",
			"Reply content",
			"Reply Description",
		}

		for _, field := range expectedFields {
			if !strings.Contains(str, field) {
				t.Errorf("String() result does not contain %q", field)
			}
		}
	})

	t.Run("ProtoMessage", func(t *testing.T) {
		req.ProtoMessage()
	})

	t.Run("ProtoReflect", func(t *testing.T) {
		r := req.ProtoReflect()
		if r == nil {
			t.Error("ProtoReflect() returned nil")
		}

		if r.Descriptor().FullName() != "proto.ReplyEmailRequest" {
			t.Errorf("Wrong message name: got %v, want proto.ReplyEmailRequest",
				r.Descriptor().FullName())
		}
	})

	t.Run("Descriptor", func(t *testing.T) {
		desc, idx := req.Descriptor()
		if desc == nil {
			t.Error("Descriptor() returned nil descriptor")
		}
		if len(idx) == 0 {
			t.Error("Descriptor() returned empty index")
		}
	})
}

func TestReplyEmailReplyMethods(t *testing.T) {
	reply := &ReplyEmailReply{}

	t.Run("Reset", func(t *testing.T) {
		r := &ReplyEmailReply{}
		r.Reset()
		if !proto.Equal(r, &ReplyEmailReply{}) {
			t.Error("Reset() did not clear the message")
		}
	})

	t.Run("String", func(t *testing.T) {
		str := reply.String()
		if str != "" {
			t.Errorf("String() = %q, want empty string", str)
		}

		reply2 := &ReplyEmailReply{}
		str2 := reply2.String()
		if str2 != "" {
			t.Errorf("String() = %q, want empty string", str2)
		}
	})

	t.Run("ProtoMessage", func(t *testing.T) {
		reply.ProtoMessage()
	})

	t.Run("ProtoReflect", func(t *testing.T) {
		r := reply.ProtoReflect()
		if r == nil {
			t.Error("ProtoReflect() returned nil")
		}
		if r.Descriptor().FullName() != "proto.ReplyEmailReply" {
			t.Errorf("Wrong message name: got %v, want proto.ReplyEmailReply",
				r.Descriptor().FullName())
		}
	})

	t.Run("Descriptor", func(t *testing.T) {
		desc, idx := reply.Descriptor()
		if desc == nil {
			t.Error("Descriptor() returned nil descriptor")
		}
		if len(idx) == 0 {
			t.Error("Descriptor() returned empty index")
		}
	})
}

func TestReplyEmailRequestGetters(t *testing.T) {
	fixedTime := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	timestamp := timestamppb.New(fixedTime)

	req := &ReplyEmailRequest{
		From:        "sender@example.com",
		To:          "recipient@example.com",
		Title:       "Re: Test",
		Sender:      "original@example.com",
		ReplyText:   "Reply content",
		SendingDate: timestamp,
		Description: "Reply Description",
	}

	t.Run("GetFrom", func(t *testing.T) {
		if got := req.GetFrom(); got != "sender@example.com" {
			t.Errorf("GetFrom() = %v, want %v", got, "sender@example.com")
		}
	})

	t.Run("GetTo", func(t *testing.T) {
		if got := req.GetTo(); got != "recipient@example.com" {
			t.Errorf("GetTo() = %v, want %v", got, "recipient@example.com")
		}
	})

	t.Run("GetTitle", func(t *testing.T) {
		if got := req.GetTitle(); got != "Re: Test" {
			t.Errorf("GetTitle() = %v, want %v", got, "Re: Test")
		}
	})

	t.Run("GetSender", func(t *testing.T) {
		if got := req.GetSender(); got != "original@example.com" {
			t.Errorf("GetSender() = %v, want %v", got, "original@example.com")
		}
	})

	t.Run("GetReplyText", func(t *testing.T) {
		if got := req.GetReplyText(); got != "Reply content" {
			t.Errorf("GetReplyText() = %v, want %v", got, "Reply content")
		}
	})

	t.Run("GetSendingDate", func(t *testing.T) {
		got := req.GetSendingDate()
		if got.AsTime().Unix() != fixedTime.Unix() {
			t.Errorf("GetSendingDate().Unix() = %v, want %v", got.AsTime().Unix(), fixedTime.Unix())
		}
	})

	t.Run("GetDescription", func(t *testing.T) {
		if got := req.GetDescription(); got != "Reply Description" {
			t.Errorf("GetDescription() = %v, want %v", got, "Reply Description")
		}
	})

	emptyReq := &ReplyEmailRequest{}

	t.Run("GetFromEmpty", func(t *testing.T) {
		if got := emptyReq.GetFrom(); got != "" {
			t.Errorf("GetFrom() = %v, want empty string", got)
		}
	})

	t.Run("GetToEmpty", func(t *testing.T) {
		if got := emptyReq.GetTo(); got != "" {
			t.Errorf("GetTo() = %v, want empty string", got)
		}
	})

	t.Run("GetTitleEmpty", func(t *testing.T) {
		if got := emptyReq.GetTitle(); got != "" {
			t.Errorf("GetTitle() = %v, want empty string", got)
		}
	})

	t.Run("GetSenderEmpty", func(t *testing.T) {
		if got := emptyReq.GetSender(); got != "" {
			t.Errorf("GetSender() = %v, want empty string", got)
		}
	})

	t.Run("GetReplyTextEmpty", func(t *testing.T) {
		if got := emptyReq.GetReplyText(); got != "" {
			t.Errorf("GetReplyText() = %v, want empty string", got)
		}
	})

	t.Run("GetSendingDateEmpty", func(t *testing.T) {
		if got := emptyReq.GetSendingDate(); got != nil {
			t.Errorf("GetSendingDate() = %v, want nil", got)
		}
	})

	t.Run("GetDescriptionEmpty", func(t *testing.T) {
		if got := emptyReq.GetDescription(); got != "" {
			t.Errorf("GetDescription() = %v, want empty string", got)
		}
	})
}

func TestFetchEmailsViaPOP3RequestMethods(t *testing.T) {
	req := &FetchEmailsViaPOP3Request{}

	t.Run("Reset", func(t *testing.T) {
		reqCopy := *req
		reqCopy.Reset()
		if !proto.Equal(&reqCopy, &FetchEmailsViaPOP3Request{}) {
			t.Error("Reset() did not clear the message")
		}
	})

	t.Run("String", func(t *testing.T) {
		str := req.String()
		if str != "" {
			t.Errorf("String() = %q, want empty string", str)
		}
	})

	t.Run("ProtoMessage", func(t *testing.T) {
		req.ProtoMessage()
	})

	t.Run("ProtoReflect", func(t *testing.T) {
		r := req.ProtoReflect()
		if r == nil {
			t.Error("ProtoReflect() returned nil")
		}
		if r.Descriptor().FullName() != "proto.FetchEmailsViaPOP3Request" {
			t.Errorf("Wrong message name: got %v, want proto.FetchEmailsViaPOP3Request",
				r.Descriptor().FullName())
		}
	})

	t.Run("Descriptor", func(t *testing.T) {
		desc, idx := req.Descriptor()
		if desc == nil {
			t.Error("Descriptor() returned nil descriptor")
		}
		if len(idx) == 0 {
			t.Error("Descriptor() returned empty index")
		}
	})
}

func TestFetchEmailsViaPOP3ReplyMethods(t *testing.T) {
	reply := &FetchEmailsViaPOP3Reply{}

	t.Run("Reset", func(t *testing.T) {
		replyCopy := *reply
		replyCopy.Reset()
		if !proto.Equal(&replyCopy, &FetchEmailsViaPOP3Reply{}) {
			t.Error("Reset() did not clear the message")
		}
	})

	t.Run("String", func(t *testing.T) {
		str := reply.String()
		if str != "" {
			t.Errorf("String() = %q, want empty string", str)
		}
	})

	t.Run("ProtoMessage", func(t *testing.T) {
		reply.ProtoMessage()
	})

	t.Run("ProtoReflect", func(t *testing.T) {
		r := reply.ProtoReflect()
		if r == nil {
			t.Error("ProtoReflect() returned nil")
		}
		if r.Descriptor().FullName() != "proto.FetchEmailsViaPOP3Reply" {
			t.Errorf("Wrong message name: got %v, want proto.FetchEmailsViaPOP3Reply",
				r.Descriptor().FullName())
		}
	})

	t.Run("Descriptor", func(t *testing.T) {
		desc, idx := reply.Descriptor()
		if desc == nil {
			t.Error("Descriptor() returned nil descriptor")
		}
		if len(idx) == 0 {
			t.Error("Descriptor() returned empty index")
		}
	})
}

func TestFetchEmailsViaPOP3Messages(t *testing.T) {
	t.Run("Request", func(t *testing.T) {
		original := &FetchEmailsViaPOP3Request{}

		data, err := proto.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		decoded := &FetchEmailsViaPOP3Request{}
		err = proto.Unmarshal(data, decoded)
		if err != nil {
			t.Fatalf("Failed to unmarshal request: %v", err)
		}

		if !proto.Equal(original, decoded) {
			t.Error("Decoded request differs from original")
		}
	})

	t.Run("Reply", func(t *testing.T) {
		original := &FetchEmailsViaPOP3Reply{}

		data, err := proto.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal reply: %v", err)
		}

		decoded := &FetchEmailsViaPOP3Reply{}
		err = proto.Unmarshal(data, decoded)
		if err != nil {
			t.Fatalf("Failed to unmarshal reply: %v", err)
		}

		if !proto.Equal(original, decoded) {
			t.Error("Decoded reply differs from original")
		}
	})
}
