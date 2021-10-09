package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"grpc-chat-client-v2/proto"
	"grpc-chat-client-v2/routines"
	"grpc-chat-client-v2/util"
	"log"
	"os"
	"sync"
)

// TODO: add gateway for server selection
func main() {
	done := make(chan int, 1)
	wg := new(sync.WaitGroup)

	usernameFlag := flag.String("username", "", "A username to use for the chat.")
	roomFlag := flag.String("room", "", "The room to subscribe to. Should be connected to the same gateway.")
	listRoomsFlag := flag.Bool("listRooms", false, "List available rooms on the given gateway")
	gatewayFlag := flag.String("gateway", "", "A gateway to relay your messages.")
	gatewayPortFlag := flag.String("gatewayPort", "9000", "The gateway server's port.")
	flag.Parse()
	if *usernameFlag == "" {
		if *roomFlag == "" && !*listRoomsFlag {
			flag.PrintDefaults()
			os.Exit(0)
		} else {
			if *gatewayFlag == "" {
				flag.PrintDefaults()
				os.Exit(0)
			}
		}
	}
	gatewayHost := fmt.Sprintf("%s:%s", *gatewayFlag, *gatewayPortFlag)

	conn, err := grpc.Dial(gatewayHost,
		grpc.WithInsecure(),
		//grpc.WithKeepaliveParams(keepalive.ClientParameters{
		//	Time:                5 * time.Second,
		//	Timeout:             time.Second,
		//	PermitWithoutStream: true}),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Connected to gateway at %s\n", gatewayHost)

	cl := proto.NewChatServiceClient(conn)

	if *listRoomsFlag {
		rooms, _ := cl.ListRooms(context.Background(), &proto.Empty{})
		log.Printf("Found %d rooms:\n", len(rooms.RoomIDs))
		for _, v := range rooms.RoomIDs {
			log.Printf("%+v\n", v)
		}
		return
	}
	connectionRequest := util.NewConnectionRequestForUsername(*usernameFlag)

	stream, err := cl.Subscribe(context.Background(), &proto.RoomRequest{
		RoomID:                   *roomFlag,
		InitialConnectionRequest: connectionRequest,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Joined room %s.\n\n", *roomFlag)

	go routines.ListenForMessages(done, wg, stream, *usernameFlag)
	go routines.ListenForInput(cl, wg, *usernameFlag, *roomFlag)
	wg.Wait()

	/* graceful shutdown */
	<-done
	conn.Close()
}
