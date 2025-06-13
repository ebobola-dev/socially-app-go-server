package config

import (
	"fmt"
	"strings"
)

type BuildType int

const (
	Development BuildType = iota
	Test
	Production
)

func ParseBuildType(s string) (BuildType, error) {
	switch strings.ToLower(s) {
	case "dev":
		return Development, nil
	case "test":
		return Test, nil
	case "prod":
		return Production, nil
	default:
		return Development, fmt.Errorf("invalid build type: %s", s)
	}
}
