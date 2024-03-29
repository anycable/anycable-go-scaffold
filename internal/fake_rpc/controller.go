package fake_rpc

import (
	"encoding/json"

	"github.com/anycable/anycable-go/common"
	"github.com/anycable/anycable-go/node"
	"github.com/anycable/anycable-go/utils"
	"github.com/apex/log"
)

const (
	welcomeMessage = "{\"type\":\"welcome\"}"
)

type Controller struct {
	log *log.Entry
}

var _ node.Controller = (*Controller)(nil)

func NewController() *Controller {
	return &Controller{log: log.WithField("context", "fake_rpc")}
}

// Start is no-op
func (c *Controller) Start() error {
	c.log.Warnf("Using fake RPC controller")
	return nil
}

// Shutdown is no-op
func (c *Controller) Shutdown() error {
	return nil
}

func (c *Controller) Authenticate(sid string, env *common.SessionEnv) (*common.ConnectResult, error) {
	c.log.WithField("sid", sid).Debug("> Connected")

	return &common.ConnectResult{
		Status:        common.SUCCESS,
		Identifier:    sid,
		Transmissions: []string{welcomeMessage},
	}, nil
}

func (c *Controller) Subscribe(sid string, env *common.SessionEnv, identifiers string, channel string) (*common.CommandResult, error) {
	c.log.WithField("sid", sid).Debugf("> Subscribed to %s", channel)

	res := &common.CommandResult{
		Status:        common.SUCCESS,
		Transmissions: []string{cableMessage("confirm_subscription", channel)},
	}
	return res, nil
}

func (c *Controller) Unsubscribe(sid string, env *common.SessionEnv, identifiers string, channel string) (*common.CommandResult, error) {
	c.log.WithField("sid", sid).Debugf("> Unubscribed from %s", channel)

	res := &common.CommandResult{
		Status: common.SUCCESS,
	}
	return res, nil
}

func (c *Controller) Perform(sid string, env *common.SessionEnv, id string, channel string, data string) (res *common.CommandResult, err error) {
	var payload map[string]interface{}

	if err = json.Unmarshal([]byte(data), &payload); err != nil {
		return nil, err
	}

	action := payload["action"].(string)

	c.log.WithField("sid", sid).Debugf("> Perform action: %s, data: %v", action, payload)

	nextState := make(map[string]string)

	res = &common.CommandResult{
		Status:         common.SUCCESS,
		Disconnect:     false,
		StopAllStreams: false,
		Streams:        nil,
		Transmissions:  []string{},
		IState:         nextState,
	}

	return res, nil
}

func (c *Controller) Disconnect(sid string, env *common.SessionEnv, id string, subscriptions []string) error {
	c.log.WithField("sid", sid).Debug("> Disconnected")

	return nil
}

func cableMessage(typ string, identifier string) string {
	msg := common.Reply{Identifier: identifier, Type: typ}

	return string(utils.ToJSON(msg))
}
