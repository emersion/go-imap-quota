package quota

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// Client is a QUOTA client.
type Client struct {
	c *client.Client
}

// NewClient creates a new client.
func NewClient(c *client.Client) *Client {
	return &Client{c: c}
}

// SupportQuota checks if the server supports the QUOTA extension.
func (c *Client) SupportQuota() (bool, error) {
	return c.c.Support(Capability)
}

// SetQuota changes the resource limits for the specified quota root. Any
// previous resource limits for the named quota root are discarded.
func (c *Client) SetQuota(root string, resources map[string]uint32) error {
	if c.c.State&imap.AuthenticatedState == 0 {
		return client.ErrNotLoggedIn
	}

	cmd := &SetCommand{
		Root:      root,
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
	if c.c.State&imap.AuthenticatedState == 0 {
		return nil, client.ErrNotLoggedIn
	}

	cmd := &GetCommand{
		Root: root,
	}

	res := &Response{}

	status, err := c.c.Execute(cmd, res)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}
	if len(res.Quotas) != 1 {
		return nil, fmt.Errorf("Expected exactly one QUOTA response, got %v", len(res.Quotas))
	}

	return res.Quotas[0], nil
}

// GetQuotaRoot returns the list of quota roots for a mailbox.
func (c *Client) GetQuotaRoot(mailbox string) ([]*Status, error) {
	if c.c.State&imap.AuthenticatedState == 0 {
		return nil, client.ErrNotLoggedIn
	}

	cmd := &GetRootCommand{
		Mailbox: mailbox,
	}

	res := &Response{}

	status, err := c.c.Execute(cmd, res)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}

	return res.Quotas, nil
}
