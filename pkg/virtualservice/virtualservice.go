package virtualservice

const (
	DefaultVirtualService = "default"
)

// Options represents the CLI parameters for virtual services
type Options struct {
	Filename    string
	Output      string
	Template    string
	Interactive bool
}
