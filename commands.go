package quota

import (
	"errors"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/utf7"
)

const (
	setCommandName     = "SETQUOTA"
	getCommandName     = "GETQUOTA"
	getRootCommandName = "GETQUOTAROOT"
)

// The SETQUOTA command. See RFC 2087 section 4.1.
type SetCommand struct {
	Root      string
	Resources map[string]uint32
}

func (cmd *SetCommand) Command() *imap.Command {
	args := []interface{}{cmd.Root}

	for k, v := range cmd.Resources {
		args = append(args, k, v)
	}

	return &imap.Command{
		Name:      setCommandName,
		Arguments: args,
	}
}

func (cmd *SetCommand) Parse(fields []interface{}) error {
	if len(fields) < 2 {
		return errors.New("No enough arguments")
	}

	var ok bool
	if cmd.Root, ok = fields[0].(string); !ok {
		return errors.New("Quota root must be a string")
	}

	resources, ok := fields[1].([]interface{})
	if !ok {
		return errors.New("Resources must be a list")
	}

	var name string
	for i, v := range resources {
		if i%2 == 0 {
			name, ok = v.(string)
			if !ok {
				return errors.New("Resource name must be a string")
			}
		} else {
			var err error
			cmd.Resources[name], err = imap.ParseNumber(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// The GETQUOTA command. See RFC 2087 section 4.2.
type GetCommand struct {
	Root string
}

func (cmd *GetCommand) Command() *imap.Command {
	return &imap.Command{
		Name:      getCommandName,
		Arguments: []interface{}{cmd.Root},
	}
}

func (cmd *GetCommand) Parse(fields []interface{}) error {
	if len(fields) < 1 {
		return errors.New("No enough arguments")
	}

	var ok bool
	if cmd.Root, ok = fields[0].(string); !ok {
		return errors.New("Quota root must be a string")
	}

	return nil
}

// The GETQUOTAROOT command. See RFC 2087 section 4.3.
type GetRootCommand struct {
	Mailbox string
}

func (cmd *GetRootCommand) Command() *imap.Command {
	mailbox, _ := utf7.Encoder.String(cmd.Mailbox)

	return &imap.Command{
		Name:      getRootCommandName,
		Arguments: []interface{}{mailbox},
	}
}

func (cmd *GetRootCommand) Parse(fields []interface{}) error {
	if len(fields) < 1 {
		return errors.New("No enough arguments")
	}

	var ok bool
	mailbox, ok := fields[0].(string)
	if !ok {
		return errors.New("Quota root must be a string")
	}
	var err error
	if cmd.Mailbox, err = utf7.Decoder.String(mailbox); err != nil {
		return err
	}

	return nil
}
