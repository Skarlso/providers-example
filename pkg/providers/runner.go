package providers

// Runner runs a plugin.
type Runner interface {
	Run(args []string) error
}
