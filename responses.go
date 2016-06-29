package quota

import (
	"errors"

	"github.com/emersion/go-imap/common"
	"github.com/emersion/go-imap/utf7"
)

const (
	responseName = "QUOTA"
	rootResponseName = "QUOTAROOT"
)

// A quota status.
type Status struct {
	// The quota name.
	Name string
	// The quota resources. Each resource is indexed by its name and contains its
	// current usage as well as its limit.
	Resources map[string][2]uint32
}

func (rs *Status) Parse(fields []interface{}) error {
	if len(fields) < 2 {
		return errors.New("No enough arguments")
	}

	var ok bool
	if rs.Name, ok = fields[0].(string); !ok {
		return errors.New("Quota root must be a string")
	}

	resources, ok := fields[1].([]interface{})
	if !ok {
		return errors.New("Resources must be a list")
	}

	var name string
	var usage, limit uint32
	var err error
	for i, v := range resources {
		if ii := i % 3; ii == 0 {
			name, ok = v.(string)
			if !ok {
				return errors.New("Resource name must be a string")
			}
		} else if ii == 1 {
			if usage, err = common.ParseNumber(v); err != nil {
				return err
			}
		} else {
			if limit, err = common.ParseNumber(v); err != nil {
				return err
			}
			rs.Resources[name] = [2]uint32{usage, limit}
		}
	}

	return nil
}

func (rs *Status) Format() (fields []interface{}) {
	fields = append(fields, rs.Name)
	for k, v := range rs.Resources {
		fields = append(fields, k, v[0], v[1])
	}
	return
}

// A QUOTA response. See RFC 2087 section 5.1.
type Response struct {
	Quotas chan *Status
}

func (r *Response) HandleFrom(hdlr common.RespHandler) (err error) {
	for h := range hdlr {
		fields, ok := h.AcceptNamedResp(responseName)
		if !ok {
			continue
		}

		quota := &Status{}
		if err = quota.Parse(fields); err != nil {
			return
		}

		r.Quotas <- quota
	}

	return
}

func (r *Response) WriteTo(w *common.Writer) (err error) {
	for quota := range r.Quotas {
		fields := []interface{}{responseName}
		fields = append(fields, quota.Format()...)

		res := common.NewUntaggedResp(fields)
		if err = res.WriteTo(w); err != nil {
			return
		}
	}

	return
}

type MailboxRoots struct {
	Name string
	Roots []string
}

func (m *MailboxRoots) Parse(fields []interface{}) error {
	if len(fields) < 1 {
		return errors.New("No enough arguments")
	}

	mailbox, ok := fields[0].(string)
	if !ok {
		return errors.New("Mailbox name must be a string")
	}
	var err error
	if m.Name, err = utf7.Decoder.String(mailbox); err != nil {
		return err
	}

	for _, f := range fields[1:] {
		root, ok := f.(string)
		if !ok {
			return errors.New("Quota root must be a string")
		}
		m.Roots = append(m.Roots, root)
	}

	return nil
}

func (m *MailboxRoots) Format() (fields []interface{}) {
	fields = append(fields, m.Name)
	for _, root := range m.Roots {
		fields = append(fields, root)
	}
	return
}

// A QUOTAROOT response. See RFC 2087 section 5.1.
type RootResponse struct {
	Mailbox *MailboxRoots
}

func (r *RootResponse) HandleFrom(hdlr common.RespHandler) error {
	for h := range hdlr {
		fields, ok := h.AcceptNamedResp(rootResponseName)
		if !ok {
			continue
		}

		m := &MailboxRoots{}
		if err := m.Parse(fields); err != nil {
			return err
		}

		r.Mailbox = m
	}

	return nil
}

func (r *RootResponse) WriteTo(w *common.Writer) (err error) {
	fields := []interface{}{rootResponseName}
	fields = append(fields, r.Mailbox.Format()...)

	res := common.NewUntaggedResp(fields)
	if err = res.WriteTo(w); err != nil {
		return
	}

	return
}
