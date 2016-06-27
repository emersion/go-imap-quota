// Implements the IMAP QUOTA extension, as defined in RFC 2087.
package quota

// The QUOTA capability name.
const Capability = "QUOTA"

// Resources defined in RFC 2087 section 3.
const (
	// Sum of messages' RFC822.SIZE, in units of 1024 octets
	ResourceStorage = "STORAGE"
	// Number of messages
	ResourceMessage = "MESSAGE"
)
