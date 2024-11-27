package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginRequest(t *testing.T) {
	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	t.Run("GetEmail", func(t *testing.T) {
		assert.Equal(t, "test@example.com", req.GetEmail())
	})

	t.Run("GetPassword", func(t *testing.T) {
		assert.Equal(t, "password123", req.GetPassword())
	})

	t.Run("Reset", func(t *testing.T) {
		req.Reset()
		assert.Empty(t, req.Email)
		assert.Empty(t, req.Password)
	})
}

func TestLoginReply(t *testing.T) {
	reply := &LoginReply{
		Avatar:    "avatar.jpg",
		SessionId: "session123",
		CsrfId:    "csrf123",
		Name:      "Test User",
	}

	t.Run("GetAvatar", func(t *testing.T) {
		assert.Equal(t, "avatar.jpg", reply.GetAvatar())
	})

	t.Run("GetSessionId", func(t *testing.T) {
		assert.Equal(t, "session123", reply.GetSessionId())
	})

	t.Run("GetCsrfId", func(t *testing.T) {
		assert.Equal(t, "csrf123", reply.GetCsrfId())
	})

	t.Run("GetName", func(t *testing.T) {
		assert.Equal(t, "Test User", reply.GetName())
	})

	t.Run("Reset", func(t *testing.T) {
		reply.Reset()
		assert.Empty(t, reply.Avatar)
		assert.Empty(t, reply.SessionId)
		assert.Empty(t, reply.CsrfId)
		assert.Empty(t, reply.Name)
	})
}

func TestString(t *testing.T) {
	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	assert.NotEmpty(t, req.String())

	reply := &LoginReply{
		Avatar:    "avatar.jpg",
		SessionId: "session123",
		CsrfId:    "csrf123",
		Name:      "Test User",
	}
	assert.NotEmpty(t, reply.String())
}

func TestProtoMessage(t *testing.T) {
	req := &LoginRequest{}
	req.ProtoMessage()

	reply := &LoginReply{}
	reply.ProtoMessage()
}

func TestProtoReflect(t *testing.T) {
	req := &LoginRequest{}
	assert.NotNil(t, req.ProtoReflect())

	reply := &LoginReply{}
	assert.NotNil(t, reply.ProtoReflect())
}

func TestDescriptor(t *testing.T) {
	assert.NotNil(t, File_auth_proto)
}
