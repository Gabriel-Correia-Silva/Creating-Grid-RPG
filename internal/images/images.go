package images

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

func Load(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return png.Decode(f)
	case ".jpg", ".jpeg":
		return jpeg.Decode(f)
	case ".bmp":
		return bmp.Decode(f)
	case ".gif":
		g, err := gif.DecodeAll(f)
		if err != nil {
			return nil, err
		}
		return g.Image[0], nil
	case ".tif", ".tiff":
		return tiff.Decode(f)
	case ".webp":
		return webp.Decode(f)
	default:
		img, _, err := image.Decode(f)
		return img, err
	}
}

func Size(path string) (int, int, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".gif", ".webp", ".tif", ".tiff":
		img, err := Load(path)
		if err != nil {
			return 0, 0, err
		}
		b := img.Bounds()
		return b.Dx(), b.Dy(), nil
	}

	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

func ConvertToPNGTemp(path string) (string, error) {
	img, err := Load(path)
	if err != nil {
		return "", fmt.Errorf("erro ao carregar imagem: %w", err)
	}

	tmp, err := os.CreateTemp("", "grid-bg-*.png")
	if err != nil {
		return "", fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer tmp.Close()

	if err := png.Encode(tmp, img); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("erro ao codificar PNG temporário: %w", err)
	}
	return tmp.Name(), nil
}

func IsPDFNative(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
