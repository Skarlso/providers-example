package models

const (
	// Bare metal plugin type.
	Bare = "bare"
	// Container plugin type.
	Container = "container"
)

// Plugin defines what a Plugin looks like.
type Plugin struct {
	ID        int
	Name      string
	Type      string
	Container *ContainerPlugin
	Bare      *BareMetalPlugin
}

// ContainerPlugin is a specific plugin which is in a container.
type ContainerPlugin struct {
	Image string
}

// BareMetalPlugin is a plugin which is a file on the filesystem.
type BareMetalPlugin struct {
	Location string
}
