# go-imap-quota

[![GoDoc](https://godoc.org/github.com/emersion/go-imap-quota?status.svg)](https://godoc.org/github.com/emersion/go-imap-quota)

[QUOTA extension](https://tools.ietf.org/html/rfc2087) for [go-imap](https://github.com/emersion/go-imap)

## Usage

```go
package main

import (
	"log"

	"github.com/emersion/go-imap-quota"
)

func main() {
	// Connect to IMAP server

	// Create a quota client
	qc := quota.NewClient(c)

	// Check for server support
	if !qc.SupportsQuota() {
		log.Fatal("Client doesn't support QUOTA extension")
	}

	// Retrieve quotas for INBOX
	quotas, err := qc.GetQuotaRoot("INBOX")
	if err != nil {
		log.Fatal(err)
	}

	// Print quotas
	log.Println("Quotas for INBOX:")
	for _, quota := range quotas {
		log.Printf("* %q, resources:\n", quota.Name)
		for name, usage := range quota.Resources {
			log.Printf("  * %v: %v/%v used\n", name, usage[0], usage[1])
		}
	}
}
```

## License

MIT
