package domain

import "time"

type Midia struct {
	IDMidia     string // The 'abc123def' part of the mxc:// URI
	Origin      string // Your server name, or a remote server name if federated
	ContentType string // e.g., "image/jpeg"
	SizeBytes   int64
	UploadName  string // Original file name
	IDUsuario   string // Who uploaded it
	CreatedAt   time.Time
}
