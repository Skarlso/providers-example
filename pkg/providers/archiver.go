package providers

// Archiver can extract files from an archive.
type Archiver interface {
	Untar(content []byte) ([]byte, error)
}
