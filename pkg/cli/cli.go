package cli

import (
	"net/http"

	acli "github.com/anycable/anycable-go/cli"
	aconfig "github.com/anycable/anycable-go/config"
	"github.com/anycable/anycable-go/metrics"
	"github.com/anycable/anycable-go/node"
	"github.com/anycable/anycable-go/server"
	"github.com/anycable/anycable-go/ws"
	"github.com/apex/log"
	"github.com/gorilla/websocket"

	"github.com/anycable/mycable/internal/fake_rpc"
	"github.com/anycable/mycable/pkg/config"
	"github.com/anycable/mycable/pkg/custom"
	"github.com/anycable/mycable/pkg/version"
)

func Run(conf *config.Config, anyconf *aconfig.Config) error {
	anycable, err := initAnyCableRunner(conf, anyconf)

	if err != nil {
		return err
	}

	log.WithField("context", "main").Infof("Starting custom AnyCable v%s", version.Version())

	return anycable.Run()
}

func initAnyCableRunner(appConf *config.Config, anyConf *aconfig.Config) (*acli.Runner, error) {
	opts := []acli.Option{
		acli.WithName("AnyCable"),
		acli.WithDefaultSubscriber(),
		acli.WithDefaultBroker(),
		// Enable broadcasting
		// acli.WithDefaultBroadcaster(),
		acli.WithWebSocketEndpoint("/ws", myWebsocketHandler(appConf)),
	}

	if appConf.FakeRPC {
		opts = append(opts, acli.WithController(func(m *metrics.Metrics, c *aconfig.Config) (node.Controller, error) {
			return fake_rpc.NewController(), nil
		}))
	} else {
		opts = append(opts, acli.WithDefaultRPCController())
	}

	return acli.NewRunner(anyConf, opts)
}

func myWebsocketHandler(config *config.Config) func(n *node.Node, c *aconfig.Config) (http.Handler, error) {
	return func(n *node.Node, c *aconfig.Config) (http.Handler, error) {
		extractor := server.DefaultHeadersExtractor{Headers: c.Headers, Cookies: c.Cookies}

		executor := custom.NewExecutor(n)

		log.WithField("context", "main").Infof("Handle custom WebSocket connections at ws://%s:%d/ws", c.Host, c.Port)

		return ws.WebsocketHandler([]string{}, &extractor, &c.WS, func(wsc *websocket.Conn, info *server.RequestInfo, callback func()) error {
			wrappedConn := ws.NewConnection(wsc)
			session := node.NewSession(
				n, wrappedConn, info.URL, info.Headers, info.UID,
				node.WithEncoder(custom.Encoder{}), node.WithExecutor(executor),
			)

			_, err := n.Authenticate(session)

			if err != nil {
				return err
			}

			return session.Serve(callback)
		}), nil
	}
}
