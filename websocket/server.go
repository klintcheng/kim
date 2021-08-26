package websocket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/segmentio/ksuid"
)

// ServerOptions ServerOptions
type ServerOptions struct {
	loginwait time.Duration //登录超时
	readwait  time.Duration //读超时
	writewait time.Duration //写超时
}

// Server is a websocket implement of the Server
type Server struct {
	listen string
	kim.ServiceRegistration
	kim.ChannelMap
	kim.Acceptor
	kim.MessageListener
	kim.StateListener
	once    sync.Once
	options ServerOptions
}

// NewServer NewServer
func NewServer(listen string, service kim.ServiceRegistration) kim.Server {
	return &Server{
		listen:              listen,
		ServiceRegistration: service,
		options: ServerOptions{
			loginwait: kim.DefaultLoginWait,
			readwait:  kim.DefaultReadWait,
			writewait: kim.DefaultWriteWait,
		},
	}
}

// Start server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	log := logger.WithFields(logger.Fields{
		"module": "ws.server",
		"listen": s.listen,
		"id":     s.ServiceID(),
	})

	if s.Acceptor == nil {
		s.Acceptor = new(defaultAcceptor)
	}
	if s.StateListener == nil {
		return fmt.Errorf("StateListener is nil")
	}
	if s.ChannelMap == nil {
		s.ChannelMap = kim.NewChannels(100)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// step 1
		rawconn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			resp(w, http.StatusBadRequest, err.Error())
			return
		}

		// step 2 包装conn
		conn := NewConn(rawconn)
		// step 3
		id, err := s.Accept(conn, s.options.loginwait)
		if err != nil {
			log.Warn(err)
			_ = conn.WriteFrame(kim.OpClose, []byte(err.Error()))
			conn.Close()
			return
		}
		if _, ok := s.Get(id); ok {
			log.Warnf("channel %s existed", id)
			_ = conn.WriteFrame(kim.OpClose, []byte("channelId is repeated"))
			conn.Close()
			return
		}
		// step 4
		channel := kim.NewChannel(id, conn)
		channel.SetWriteWait(s.options.writewait)
		channel.SetReadWait(s.options.readwait)
		s.Add(channel)

		go func(ch kim.Channel) {
			// step 5
			err := ch.Readloop(s.MessageListener)
			if err != nil {
				log.Info(err)
			}
			// step 6
			s.Remove(ch.ID())
			err = s.Disconnect(ch.ID())
			if err != nil {
				log.Warn(err)
			}
			ch.Close()
		}(channel)
	})
	log.Infoln("started")
	return http.ListenAndServe(s.listen, mux)
}

// Shutdown Shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	log := logger.WithFields(logger.Fields{
		"module": "ws.server",
		"id":     s.ServiceID(),
	})
	s.once.Do(func() {
		defer func() {
			log.Infoln("shutdown")
		}()
		// close channels
		chanels := s.ChannelMap.All()
		for _, ch := range chanels {
			ch.Close()

			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}

	})
	return nil
}

// string channelID
// []byte data
func (s *Server) Push(id string, data []byte) error {
	ch, ok := s.ChannelMap.Get(id)
	if !ok {
		return errors.New("channel no found")
	}
	return ch.Push(data)
}

// SetAcceptor SetAcceptor
func (s *Server) SetAcceptor(acceptor kim.Acceptor) {
	s.Acceptor = acceptor
}

// SetMessageListener SetMessageListener
func (s *Server) SetMessageListener(listener kim.MessageListener) {
	s.MessageListener = listener
}

// SetStateListener SetStateListener
func (s *Server) SetStateListener(listener kim.StateListener) {
	s.StateListener = listener
}

// SetChannels SetChannels
func (s *Server) SetChannelMap(channels kim.ChannelMap) {
	s.ChannelMap = channels
}

// SetReadWait set read wait duration
func (s *Server) SetReadWait(readwait time.Duration) {
	s.options.readwait = readwait
}

func resp(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	if body != "" {
		_, _ = w.Write([]byte(body))
	}
	logger.Warnf("response with code:%d %s", code, body)
}

type defaultAcceptor struct {
}

// Accept defaultAcceptor
func (a *defaultAcceptor) Accept(conn kim.Conn, timeout time.Duration) (string, error) {
	return ksuid.New().String(), nil
}
