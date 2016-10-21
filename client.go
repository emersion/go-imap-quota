package quota

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Client struct {
	c *client.Client
}

// NewClient creates a new client.
func NewClient(c *client.Client) *Client {
	return &Client{c: c}
}

// SupportsMove checks if the server supports the QUOTA extension.
func (c *Client) SupportsQuota() bool {
	return c.c.Caps[Capability]
}

// SetQuota changes the resource limits for the specified quota root. Any
// previous resource limits for the named quota root are discarded.
func (c *Client) SetQuota(root string, resources map[string]uint32) error {
	if c.c.State & imap.AuthenticatedState == 0 {
		return client.ErrNotLoggedIn
	}

	cmd := &SetCommand{
		Root: root,
		Resources: resources,
	}

	status, err := c.c.Execute(cmd, nil)
	if err != nil {
		return err
	}
	return status.Err()
}

// GetQuota returns a quota root's resource usage and limits.
func (c *Client) GetQuota(root string) (*Status, error) {
	if c.c.State & imap.AuthenticatedState == 0 {
		return nil, client.ErrNotLoggedIn
	}

	cmd := &GetCommand{
		Root: root,
	}

	ch := make(chan *Status, 1)
	res := &Response{
		Quotas: ch,
	}

	status, err := c.c.Execute(cmd, res)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}

	return <-ch, nil
}

// GetQuotaRoot returns the list of quota roots for a mailbox.
func (c *Client) GetQuotaRoot(mailbox string) (*MailboxRoots, error) {
	if c.c.State & imap.AuthenticatedState == 0 {
		return nil, client.ErrNotLoggedIn
	}

	cmd := &GetRootCommand{
		Mailbox: mailbox,
	}

	res := &RootResponse{}

	status, err := c.c.Execute(cmd, res)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}

	return res.Mailbox, nil
}
