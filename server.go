package quota

import (
	"errors"

	"github.com/emersion/go-imap"
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
	// to be the specified limits. Any previous resource limits for the named
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

func (h *SetHandler) Handle(conn server.Conn) error {
	if conn.Context().User == nil {
		return server.ErrNotAuthenticated
	}

	u, ok := conn.Context().User.(User)
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

func (h *GetHandler) Handle(conn server.Conn) error {
	if conn.Context().User == nil {
		return server.ErrNotAuthenticated
	}

	u, ok := conn.Context().User.(User)
	if !ok {
		return ErrUnsupportedBackend
	}

	status, err := u.GetQuota(h.Root)
	if err != nil {
		return err
	}

	res := &Response{Quotas: []*Status{status}}
	return conn.WriteResp(res)
}

type GetRootHandler struct {
	GetRootCommand
}

func (h *GetRootHandler) Handle(conn server.Conn) error {
	if conn.Context().User == nil {
		return server.ErrNotAuthenticated
	}

	u, ok := conn.Context().User.(User)
	if !ok {
		return ErrUnsupportedBackend
	}

	mbox, err := conn.Context().User.GetMailbox(h.Mailbox)
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
			Name:  h.Mailbox,
			Roots: roots,
		},
	}
	if err := conn.WriteResp(rootRes); err != nil {
		return err
	}

	res := &Response{}
	for _, name := range roots {
		status, err := u.GetQuota(name)
		if err != nil {
			return err
		}
		res.Quotas = append(res.Quotas, status)
	}

	return conn.WriteResp(res)
}

type extension struct{}

func NewExtension() server.Extension {
	return &extension{}
}

func (ext *extension) Capabilities(c server.Conn) []string {
	if c.Context().State&imap.AuthenticatedState != 0 {
		return []string{Capability}
	}
	return nil
}

func (ext *extension) Command(name string) server.HandlerFactory {
	switch name {
	case setCommandName:
		return func() server.Handler {
			return &SetHandler{}
		}
	case getCommandName:
		return func() server.Handler {
			return &GetHandler{}
		}
	case getRootCommandName:
		return func() server.Handler {
			return &GetRootHandler{}
		}
	}

	return nil
}
