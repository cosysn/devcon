package builder

import (
	"fmt"
)

func NewBuilder(provider string) (Builder, error) {
	switch provider {
	case "docker":
		return NewDockerBuilder()
	// case "buildkit":
	//     return NewBuildKitBuilder()
	// case "kaniko":
	//     return NewKanikoBuilder()
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}
