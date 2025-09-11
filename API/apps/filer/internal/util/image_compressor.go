package util

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/nfnt/resize"
)

type CompressionConfig struct {
	MaxWidth      int
	MaxHeight     int
	MaxFileSizeKB int
	Quality       int
}

var DefaultCompressionConfig = CompressionConfig{
	MaxWidth:      1920,
	MaxHeight:     1080,
	MaxFileSizeKB: 500,
	Quality:       85,
}

func CompressImage(data []byte, config CompressionConfig) ([]byte, string, error) {
	contentType := http.DetectContentType(data)

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	needsResize := width > config.MaxWidth || height > config.MaxHeight
	needsCompression := len(data) > config.MaxFileSizeKB*1024

	if !needsResize && !needsCompression {
		return data, contentType, nil
	}

	var resizedImage image.Image = img

	if needsResize {
		if width > height {
			resizedImage = resize.Resize(uint(config.MaxWidth), 0, img, resize.Lanczos3)
		} else {
			resizedImage = resize.Resize(0, uint(config.MaxHeight), img, resize.Lanczos3)
		}
	}

	var compressedData bytes.Buffer
	var newContentType string

	switch contentType {
	case "image/jpeg", "image/jpg":
		newContentType = "image/jpeg"
		err = jpeg.Encode(&compressedData, resizedImage, &jpeg.Options{
			Quality: config.Quality,
		})

	case "image/png":
		newContentType = "image/png"
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(&compressedData, resizedImage)

	default:
		return data, contentType, nil
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to compress image: %w", err)
	}

	compressedBytes := compressedData.Bytes()
	if len(compressedBytes) > config.MaxFileSizeKB*1024 && contentType == "image/jpeg" {
		for quality := config.Quality - 10; quality >= 50; quality -= 10 {
			var aggressiveBuffer bytes.Buffer
			err = jpeg.Encode(&aggressiveBuffer, resizedImage, &jpeg.Options{
				Quality: quality,
			})
			if err == nil && aggressiveBuffer.Len() <= config.MaxFileSizeKB*1024 {
				compressedBytes = aggressiveBuffer.Bytes()
				break
			}
		}
	}

	return compressedBytes, newContentType, nil
}
