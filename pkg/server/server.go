package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/enix/tsigan/pkg/adapters"
	"github.com/enix/tsigan/pkg/tsig"
	miekgdns "github.com/miekg/dns"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger // TODO move to Server struct

// FIXME refactor server state
type Server struct {
	Configuration  *Configuration
	keyring        tsig.TsigKeyring
	defaultKeyName string
	adapters       []*adapters.IAdapter
	adaptersByName map[string]*adapters.IAdapter
	defaultAdapter *adapters.IAdapter
	zones          []*Zone
	zonesByFqdn    map[string]*Zone
}

func NewServer(configuration *Configuration) *Server {
	return &Server{
		Configuration:  configuration,
		keyring:        tsig.NewTsigKeyring(),
		adaptersByName: make(map[string]*adapters.IAdapter),
		zonesByFqdn:    make(map[string]*Zone),
	}
}

func (s *Server) Run() error {
	if err := s.init(); err != nil {
		Logger.Fatalf("failed to init server: %s", err)
	}

	// FIXME add logging here
	tsigProvider := NewTsigProvider(&s.keyring)

	// FIXME add logging here
	miekgdns.HandleFunc(".",
		func(w miekgdns.ResponseWriter, r *miekgdns.Msg) {
			Query{
				server:   s,
				writer:   w,
				received: r,
			}.Handle()
		})

	// FIXME add logging here
	go s.serve("udp", tsigProvider, true)
	go s.serve("tcp", tsigProvider, true)

	sig := make(chan os.Signal)
	// FIXME sigchanyzer: misuse of unbuffered os.Signal channel as argument to signal.Notify
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	signal := <-sig
	Logger.Warnw("received stop signal while running server", "signal", signal)

	return nil
}

func (s *Server) serve(net string, provider *TsigProvider, soreuseport bool) {
	server := &miekgdns.Server{
		Addr:          "[::]:5353", // FIXME config from cobra and ServerSettings
		Net:           net,
		ReusePort:     soreuseport,
		TsigProvider:  provider,
		MsgAcceptFunc: s.msgAcceptAction,
		MsgInvalidFunc: func(m []byte, err error) {
			Logger.Debugw("an invalid message was observed",
				"protocol", net, "length", len(m), "error", err.Error())
		},
	}

	// FIXME add logging here
	if err := server.ListenAndServe(); err != nil {
		Logger.Fatalw("failed to start network server", "protocol", net, "error", err)
	}
}

func (s *Server) msgAcceptAction(dh miekgdns.Header) miekgdns.MsgAcceptAction {
	const (
		// from: https://github.com/miekg/dns/blob/master/types.go
		// Header.Bits
		_QR = 1 << 15 // query/response (response=1)
	)

	if isResponse := dh.Bits&_QR != 0; isResponse {
		return miekgdns.MsgIgnore
	}

	// only accept DNS updates
	opcode := int(dh.Bits>>11) & 0xF
	if opcode == miekgdns.OpcodeUpdate {
		return miekgdns.MsgAccept
	}

	return miekgdns.MsgReject
}
