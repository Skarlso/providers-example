package models

type Type string

var (
	// Bare metal plugin type.
	Bare Type = "bare"
	// Container plugin type.
	Container Type = "container"
)

// Plugin defines what a Plugin looks like.
type Plugin struct {
	ID   int
	Name string
	Type Type
}

// ContainerPlugin is a specific plugin which is in a container.
type ContainerPlugin struct {
	Plugin
	Image string
}

// BareMetalPlugin is a plugin which is a file on the filesystem.
type BareMetalPlugin struct {
	Plugin
	Path     string
	Filename string
}
