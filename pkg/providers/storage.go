package providers

// Storer can store information about the plugins that were created.
type Storer interface {
	Create()
	Delete()
	List()
}
