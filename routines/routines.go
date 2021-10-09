package routines

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"grpc-chat-client-v2/proto"
	"grpc-chat-client-v2/util"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func ListenForInput(cl proto.ChatServiceClient, wg *sync.WaitGroup, sender, recipient string) {
	wg.Add(1)
	for {
		fmt.Printf("\n<You>:")
		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		msg := sc.Text()
		//if msg[0] == uint8('\\') {
		//	// TODO: parse user command?
		//}
		if len(msg) > 1 {
			nm := util.NewMessageFrom(sender, recipient, msg)
			_, err := cl.SendMessage(context.Background(), nm)
			if err != nil {
				log.Printf("\x1b[31mFailed sending msg -> %+v\x1b[0m\n", err)
				if strings.Contains(err.Error(), "transport: Error while dialing dial") {
					log.Println("Server shut-down. Aborting.")
				}
				break
			}
		}
	}
	wg.Done()
}

func ListenForMessages(done chan int, wg *sync.WaitGroup, messageStream proto.ChatService_ConnectClient, initiatorId string) {
	wg.Add(1)
	for {
		msg, err := messageStream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				time.Sleep(time.Second / 2)
				break
			} else {
				log.Println(err)
				break
			}
		}
		if msg != nil && msg.Timestamp != 0 && msg.SenderID != initiatorId {
			fmt.Printf("\n\x1b[32m<%s>: %s\n\x1b[0m", msg.SenderUsername, string(msg.Content))
			fmt.Printf("\n<You>:")
		}
	}
	wg.Done()
	done <- 1
}
