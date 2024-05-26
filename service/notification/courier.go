package notification

import (
	"app/api/models"
	"app/pkg/logs"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

type CourierNotifyService struct {
	upgarder websocket.Upgrader
	log      logs.LoggerInterface

	MessageQueue    chan models.OrderModel
	ConnectionsList []SocketConnection
}

func NewCourierService(log logs.LoggerInterface) *CourierNotifyService {
	return &CourierNotifyService{
		upgarder: websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: true,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		log:          log,
		MessageQueue: make(chan models.OrderModel),
	}
}

func (srv *CourierNotifyService) GetUpgrader(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := srv.upgarder.Upgrade(w, r, nil)
	if err != nil {
		srv.log.Error("could not upgrade a connection for courier", logs.Error(err))
		return nil, err
	}

	srv.ConnectionsList = append(srv.ConnectionsList, SocketConnection{
		Connection: ws,
		Type:       "courier",
	})
	return ws, nil
}

// HandleNotifications Will receive a channel that will be used for
// sending messages that will be used for sending data in real time
func (srv *CourierNotifyService) HandleNotifications() {
	for {
		select {
		case data, ok := <-srv.MessageQueue:
			if !ok {
				srv.log.Error("channel is closed (courier websocket)")
				return
			}
			for i, conn := range srv.ConnectionsList {
				if err := conn.Connection.WriteJSON(data); err != nil {
					if !errors.Is(err, websocket.ErrCloseSent) {
						srv.log.Error("could not send data to websocket connection",
							logs.Error(err))
					}

					conn.Connection.Close()
					srv.ConnectionsList = append(srv.ConnectionsList[:i], srv.ConnectionsList[i+1:]...)
				}
			}
		}
	}
}

func (srv *CourierNotifyService) WriteToQueue(model models.OrderModel) {
	select {
	case srv.MessageQueue <- model:
		srv.log.Debug("data send to channel")
	default:
		srv.log.Debug("no listener found for channel")
	}
}
