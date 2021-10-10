package util

import (
	"grpc-chat-client-v2/proto"
	"time"
)

func NewConnectionRequestForUsername(username string) *proto.ConnectionRequest {
	return &proto.ConnectionRequest{
		Username: username,
		ServerID: username,
	}
}

func NewMessageFrom(sender, forWho, content string) *proto.ChatMessage {
	return &proto.ChatMessage{
		Sender:    sender,
		Recipient: forWho,
		Content:   []byte(content),
		Timestamp: uint64(time.Now().Unix()),
	}
}
