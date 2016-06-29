package quota

import (
	"errors"

	"github.com/emersion/go-imap/common"
	"github.com/emersion/go-imap/server"
)

var (
	ErrUnsupportedBackend = errors.New("quota: backend not supported")
)

type User interface {
	// GetQuota returns the quota with the specified name.
	GetQuota(name string) (*Status, error)

	// SetQuota registers or updates a quota for this user with a set of resources
	// and their limit. The resource limits for the named quota root are changed
	// to be the specified limits.  Any previous resource limits for the named
	// quota root are discarded.
	SetQuota(name string, resources map[string]uint32) error
}

type Mailbox interface {
	// ListQuotas returns the currently active quotas for this mailbox.
	ListQuotas() ([]string, error)
}

type SetHandler struct {
	SetCommand
}

func (h *SetHandler) Handle(conn *server.Conn) error {
	if conn.User == nil {
		return server.ErrNotAuthenticated
	}

	u, ok := conn.User.(User)
	if !ok {
		return ErrUnsupportedBackend
	}

	if err := u.SetQuota(h.Root, h.Resources); err != nil {
		return err
	}

	inner := &GetHandler{}
	inner.Root = h.Root
	return inner.Handle(conn)
}

type GetHandler struct {
	GetCommand
}

func (h *GetHandler) Handle(conn *server.Conn) error {
	if conn.User == nil {
		return server.ErrNotAuthenticated
	}

	u, ok := conn.User.(User)
	if !ok {
		return ErrUnsupportedBackend
	}

	status, err := u.GetQuota(h.Root)
	if err != nil {
		return err
	}

	ch := make(chan *Status, 1)
	ch <- status
	close(ch)

	res := &Response{Quotas: ch}
	return conn.WriteResp(res)
}

type GetRootHandler struct {
	GetRootCommand
}

func (h *GetRootHandler) Handle(conn *server.Conn) error {
	if conn.User == nil {
		return server.ErrNotAuthenticated
	}

	u, ok := conn.User.(User)
	if !ok {
		return ErrUnsupportedBackend
	}

	mbox, err := conn.User.GetMailbox(h.Mailbox)
	if err != nil {
		return err
	}

	qmbox, ok := mbox.(Mailbox)
	if !ok {
		return ErrUnsupportedBackend
	}

	roots, err := qmbox.ListQuotas()
	if err != nil {
		return err
	}

	rootRes := &RootResponse{
		Mailbox: &MailboxRoots{
			Name: h.Mailbox,
			Roots: roots,
		},
	}
	if err := conn.WriteResp(rootRes); err != nil {
		return err
	}

	ch := make(chan *Status, len(roots))
	for _, name := range roots {
		status, err := u.GetQuota(name)
		if err != nil {
			return err
		}
		ch <- status
	}
	close(ch)

	res := &Response{Quotas: ch}
	return conn.WriteResp(res)
}

func NewServer(s *server.Server) {
	s.RegisterCapability(Capability, common.AuthenticatedState)

	s.RegisterCommand(setCommandName, func() server.Handler {
		return &SetHandler{}
	})
	s.RegisterCommand(getCommandName, func() server.Handler {
		return &GetHandler{}
	})
	s.RegisterCommand(getRootCommandName, func() server.Handler {
		return &GetRootHandler{}
	})
}
