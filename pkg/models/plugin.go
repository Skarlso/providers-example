package models

type Type int

var (
	// Bare metal plugin type.
	Bare Type = 0
	// Container plugin type.
	Container Type = 1
)

// Plugin defines what a Plugin looks like.
type Plugin struct {
	ID   int
	Name string
	Path string
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
