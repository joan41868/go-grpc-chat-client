package util

import (
	"grpc-chat-client-v2/proto"
	"time"
)

func NewConnectionRequestForUsername(username string) *proto.ConnectionRequest {
	return &proto.ConnectionRequest{
		User: &proto.User{
			ID:       username,
			Username: username,
		},
		ServerId: username,
	}
}

func NewMessageFrom(sender, forWho, content string) *proto.ChatMessage {
	return &proto.ChatMessage{
		SenderID:       sender,
		RecipientID:    forWho,
		Content:        []byte(content),
		Timestamp:      uint64(time.Now().Unix()),
		SenderUsername: sender,
	}
}
