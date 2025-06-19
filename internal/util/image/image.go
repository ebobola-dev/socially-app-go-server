package image_util

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"golang.org/x/image/draw"
)

var (
	ErrNotImage     = errors.New("file is not an image")
	ErrInvalidImage = errors.New("unable to decode image")
	ErrMagickFailed = errors.New("ImageMagick conversion failed")
)

type SplittedImageData struct {
	Size ImageSize
	Data []byte
}

func ValidateMime(data []byte) error {
	kind, err := filetype.Match(data)
	if err != nil || kind == filetype.Unknown || kind.MIME.Type != "image" {
		return ErrNotImage
	}
	return nil
}

func ValidateImageDecode(data []byte) error {
	_, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return ErrInvalidImage
	}
	return nil
}

func ConvertToJPEG(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ConvertWithMagick(original []byte) ([]byte, error) {
	dir := "temp/magick"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	uid := uuid.New().String()
	srcPath := filepath.Join(dir, uid+"_in.jpg")
	dstPath := filepath.Join(dir, uid+"_out.jpg")

	defer func() {
		_ = os.Remove(srcPath)
		_ = os.Remove(dstPath)
	}()

	if err := os.WriteFile(srcPath, original, 0o644); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	cmdName := "convert"
	if isWindows() {
		cmdName = "magick"
	}

	cmd := exec.Command(cmdName, srcPath+"[0]", dstPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrMagickFailed, stderr.String())
	}

	data, err := os.ReadFile(dstPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read converted file: %w", err)
	}

	return data, nil
}

func SplitImageBytes(input []byte) ([]SplittedImageData, error) {
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var result []SplittedImageData
	result = append(result, SplittedImageData{
		Size: SizeOriginal,
		Data: input,
	})

	for _, size := range AllSizesOrdered() {
		if size == SizeOriginal {
			continue
		}
		target := size.IntValue()
		if target == 0 || max(width, height) <= target {
			continue
		}

		newW, newH := resizeProportional(width, height, target)

		dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

		var buf bytes.Buffer
		err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 85})
		if err != nil {
			return nil, err
		}
		result = append(result, SplittedImageData{
			Size: size,
			Data: buf.Bytes(),
		})
	}

	return result, nil
}

func resizeProportional(w, h, target int) (int, int) {
	if h > w {
		return w * target / h, target
	}
	return target, h * target / w
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isWindows() bool {
	return os.PathSeparator == '\\'
}
