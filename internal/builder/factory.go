package builder

import (
	"fmt"
)

func NewBuilder(provider string) (Builder, error) {
	return NewBuilderWithLogger(provider, nil)
}

func NewBuilderWithLogger(provider string, logger BuildLogger) (Builder, error) {
	switch provider {
	case "docker":
		b, err := NewDockerBuilder()
		if err != nil {
			return nil, err
		}
		if logger != nil {
			b.SetLogger(logger)
		}
		return b, nil
	// case "buildkit":
	//     return NewBuildKitBuilder()
	// case "kaniko":
	//     return NewKanikoBuilder()
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}
