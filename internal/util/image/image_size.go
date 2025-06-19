package image_util

import (
	"errors"
	"fmt"
	"strings"
)

type ImageSize int

const (
	Size256 ImageSize = iota
	Size512
	Size1024
	SizeOriginal
)

func (s ImageSize) IntValue() int {
	switch s {
	case Size256:
		return 256
	case Size512:
		return 512
	case Size1024:
		return 1024
	default:
		return 0
	}
}

func (s ImageSize) String() string {
	switch s {
	case Size256:
		return "256"
	case Size512:
		return "512"
	case Size1024:
		return "1024"
	case SizeOriginal:
		return "original"
	default:
		return "unknown"
	}
}

func AllSizesOrdered() []ImageSize {
	return []ImageSize{Size256, Size512, Size1024, SizeOriginal}
}

func GetOrderedSizeFrom(requested ImageSize) []ImageSize {
	if requested == SizeOriginal {
		return []ImageSize{SizeOriginal}
	}
	all := AllSizesOrdered()
	for i, sz := range all {
		if sz == requested {
			return all[i:]
		}
	}
	return []ImageSize{SizeOriginal}
}

var ImageSizesList = []ImageSize{Size256, Size512, Size1024, SizeOriginal}

func ParseImageSize(strSize string) (ImageSize, error) {
	switch strSize {
	case "256":
		return Size256, nil
	case "512":
		return Size512, nil
	case "1024":
		return Size1024, nil
	case "original":
		return SizeOriginal, nil
	default:
		return SizeOriginal, fmt.Errorf("invalid image size: %s", strSize)
	}
}

func ParseImageSizeFallback(strSize string) ImageSize {
	switch strSize {
	case "256":
		return Size256
	case "512":
		return Size512
	case "1024":
		return Size1024
	default:
		return SizeOriginal
	}
}

func SizeFromPath(path string) (ImageSize, error) {
	splitted := strings.Split(path, "/")
	filename := splitted[len(splitted)-1]
	splittedFilename := strings.Split(filename, ".")
	if len(splittedFilename) < 2 {
		return SizeOriginal, errors.New("Unable to get size, no dot in filename")
	}
	strSize := splittedFilename[0]
	size, sizeErr := ParseImageSize(strSize)
	if sizeErr != nil {
		return SizeOriginal, sizeErr
	}
	return size, nil
}
