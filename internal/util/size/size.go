package size_util

import "fmt"

func BytesToString(size int64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
	)

	floatSize := float64(size)

	switch {
	case floatSize >= TB:
		return fmt.Sprintf("%.2f TB", floatSize/TB)
	case floatSize >= GB:
		return fmt.Sprintf("%.2f GB", floatSize/GB)
	case floatSize >= MB:
		return fmt.Sprintf("%.2f MB", floatSize/MB)
	case floatSize >= KB:
		return fmt.Sprintf("%.2f KB", floatSize/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
