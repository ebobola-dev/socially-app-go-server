package image_util

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

func GetNextAvailableSizes(requested ImageSize) []ImageSize {
	all := AllSizesOrdered()
	for i, sz := range all {
		if sz == requested {
			return all[i:]
		}
	}
	return []ImageSize{SizeOriginal}
}

var ImageSizesList = []ImageSize{Size256, Size512, Size1024, SizeOriginal}
