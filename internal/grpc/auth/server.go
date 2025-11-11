package auth

import (
	"GrpcMessangerMsgServer/pkg/logger/logger"
	"fmt"
	"github.com/google/uuid"
	chat "github.com/impeaone/GrpcChatSSO/gen/go/sso"
	"google.golang.org/grpc"
	"sync"
)

type serverAPI struct {
	chat.UnimplementedChatServer
	messages []*chat.Message
	clients  sync.Map
	chatMU   sync.Mutex
	logger   *logger.Log
}

func Register(gRPC *grpc.Server, logger *logger.Log) {
	chat.RegisterChatServer(gRPC, &serverAPI{
		logger: logger,
	})
}

func (s *serverAPI) Connect(stream chat.Chat_ConnectServer) error {
	ctx := stream.Context()
	clientID := uuid.New().String()
	authMsg, err := stream.Recv()
	if err != nil {
		s.logger.Error(fmt.Sprintf("client:%s - failed to receive auth message: %v", clientID, err), logger.GetPlace())
		return fmt.Errorf("client:%s - failed to receive auth message: %v", clientID, err)
	}
	if authMsg.Text != "JOIN" {
		s.logger.Error(fmt.Sprintf("client: %s - first message must be JOIN", clientID), logger.GetPlace())
		return fmt.Errorf("client: %s - first message must be JOIN", clientID)
	}
	username := authMsg.Username
	if username == "" {
		s.logger.Error(fmt.Sprintf("client: %s - first username must be given", clientID), logger.GetPlace())
		return fmt.Errorf("client: %s - username is required", clientID)
	}

	clientChan := make(chan *chat.Message, 500)
	s.clients.Store(clientID, clientChan)
	defer func() {
		s.clients.Delete(clientID)
		close(clientChan)
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-clientChan:
				if !ok {
					return
				}
				fmt.Println(ok)
				if err := stream.Send(msg); err != nil {
					s.logger.Error(fmt.Sprintf("client:%s - failed to send message: %v", clientID, err), logger.GetPlace())
					return
				}
			}
		}

	}()
	for {
		messg, errRecv := stream.Recv()
		if errRecv != nil {
			s.logger.Error(fmt.Sprintf("Receive message error: %v", errRecv), logger.GetPlace())
			return err
		}

		s.broadcastMessage(messg, clientChan)
	}
}

func (s *serverAPI) broadcastMessage(msg *chat.Message, SenderChan chan *chat.Message) {
	s.chatMU.Lock()
	defer s.chatMU.Unlock()

	msg.Id = uuid.New().String()

	s.messages = append(s.messages, msg)
	s.clients.Range(func(key, value interface{}) bool {
		if clientChan, ok := value.(chan *chat.Message); ok && clientChan != SenderChan {
			select {
			case clientChan <- msg:
			default:
				//TODO: хуй знает что тут, не помню уже
			}
		}
		return true
	})
}
