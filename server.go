package quota

import (
	"github.com/emersion/go-imap/server"
)

type User interface {
	// GetQuota returns the quota with the specified name.
	GetQuota(name string) (*Status, error)

	// SetQuota registers or updates a quota for this user with a set of resources
	// and their limit.
	SetQuota(name string, resources map[string]uint32) error
}

type Mailbox interface {
	// ListQuotas returns the currently active quotas for this mailbox.
	ListQuotas() ([]*Status, error)
}

type GetHandler struct {
	GetCommand
}

func (h *GetHandler) Handle(conn *server.Conn) error {
	return nil // TODO
}

type SetHandler struct {
	SetCommand
}

func (h *SetHandler) Handle(conn *server.Conn) error {
	return nil // TODO
}

type GetRootHandler struct {
	GetRootCommand
}

func (h *GetRootHandler) Handle(conn *server.Conn) error {
	return nil // TODO
}
