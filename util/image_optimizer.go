package util

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
)

// IsSupportedFormat checks if the file is a supported image format (JPG, PNG, or GIF)
func IsSupportedFormat(file multipart.File) (string, error) {
	// Check file content
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return "", err
	}
	// Reset file pointer after reading
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buf)
	if contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif" {
		return contentType, nil
	}

	return "", fmt.Errorf("unsupported file format")
}

// GetImageDimensions returns the width and height of the given image file
func GetImageDimensions(file multipart.File) (int, int, error) {
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	// Reset file pointer after reading
	if _, err := file.Seek(0, 0); err != nil {
		return 0, 0, err
	}
	return img.Width, img.Height, nil
}

// ResizeImage resizes the given image file to the specified width and height and returns it as io.ReadSeeker
func ResizeImage(file multipart.File, width int, height int, format string) (io.ReadSeeker, int64, error) {
	srcImage, err := imaging.Decode(file)
	if err != nil {
		return nil, 0, err
	}

	dstImage := imaging.Resize(srcImage, width, height, imaging.Lanczos)

	var buf bytes.Buffer
	var w io.Writer = &buf
	formatLower := strings.ToLower(format)

	switch formatLower {
	case ".jpeg", ".jpg":
		err = jpeg.Encode(w, dstImage, &jpeg.Options{Quality: 95})
	case ".png":
		err = png.Encode(w, dstImage)
	case ".gif":
		err = gif.Encode(w, dstImage, nil)
	default:
		return nil, 0, fmt.Errorf("unsupported output file format:" + format)
	}

	if err != nil {
		return nil, 0, err
	}

	return bytes.NewReader(buf.Bytes()), int64(buf.Len()), nil
}
