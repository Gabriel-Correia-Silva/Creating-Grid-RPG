package raster

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"

	"grid-generator/internal/config"
	"grid-generator/internal/images"
)

const (
	dpi96Scale = 96.0 / 72.0
	pxPerCm96  = 96.0 / 2.54
)

func Generate(cfg config.GridConfig) error {
	canvas, err := buildCanvas(cfg)
	if err != nil {
		return err
	}
	drawGrid(canvas, cfg)
	return savePNG(canvas, cfg.OutputFile)
}

func buildCanvas(cfg config.GridConfig) (*image.RGBA, error) {
	w, h, err := canvasDimensions(cfg)
	if err != nil {
		return nil, err
	}
	canvas := image.NewRGBA(image.Rect(0, 0, w, h))
	if err := fillBackground(canvas, cfg); err != nil {
		return nil, err
	}
	return canvas, nil
}

func canvasDimensions(cfg config.GridConfig) (int, int, error) {
	if cfg.BackgroundImage != "" {
		img, err := images.Load(cfg.BackgroundImage)
		if err != nil {
			return 0, 0, fmt.Errorf("erro ao carregar imagem de fundo: %w", err)
		}
		b := img.Bounds()
		return b.Dx(), b.Dy(), nil
	}
	return int(math.Round(cfg.PageWidth * dpi96Scale)), int(math.Round(cfg.PageHeight * dpi96Scale)), nil
}

func fillBackground(canvas *image.RGBA, cfg config.GridConfig) error {
	if cfg.BackgroundImage != "" {
		bg, err := images.Load(cfg.BackgroundImage)
		if err != nil {
			return fmt.Errorf("erro ao carregar imagem de fundo: %w", err)
		}
		draw.Draw(canvas, canvas.Bounds(), bg, bg.Bounds().Min, draw.Src)
		return nil
	}

	var bg color.Color = color.White
	if cfg.BackgroundColor != nil {
		bg = color.RGBA{
			R: uint8(cfg.BackgroundColor.R),
			G: uint8(cfg.BackgroundColor.G),
			B: uint8(cfg.BackgroundColor.B),
			A: 255,
		}
	}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)
	return nil
}

func drawGrid(canvas *image.RGBA, cfg config.GridConfig) {
	bounds := canvas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	pxPerCm := resolvePixelsPerCm(cfg, w)
	spacingPx := cfg.GridCellCm * pxPerCm
	marginPx := cfg.MarginCm * pxPerCm
	thickness := int(math.Max(1, math.Round(cfg.LineWidth)))
	c := toRGBA(cfg.LineColor, cfg.LineOpacity)

	yStart := int(math.Round(marginPx))
	yEnd := int(math.Round(float64(h) - marginPx))
	xStart := int(math.Round(marginPx))
	xEnd := int(math.Round(float64(w) - marginPx))

	for x := marginPx; x <= float64(w)-marginPx; x += spacingPx {
		drawVerticalLine(canvas, int(math.Round(x)), yStart, yEnd, thickness, c)
	}
	for y := marginPx; y <= float64(h)-marginPx; y += spacingPx {
		drawHorizontalLine(canvas, int(math.Round(y)), xStart, xEnd, thickness, c)
	}
}

func resolvePixelsPerCm(cfg config.GridConfig, canvasWidth int) float64 {
	if cfg.BackgroundImage != "" {
		return pxPerCm96
	}
	return float64(canvasWidth) / (cfg.PageWidth / config.CmToPoints)
}

func toRGBA(col config.RGB, opacity float64) color.RGBA {
	return color.RGBA{
		R: uint8(col.R),
		G: uint8(col.G),
		B: uint8(col.B),
		A: uint8(math.Round(opacity * 255)),
	}
}

func drawVerticalLine(img *image.RGBA, x, yStart, yEnd, thickness int, c color.RGBA) {
	half := thickness / 2
	bounds := img.Bounds()
	for t := -half; t <= half; t++ {
		px := x + t
		if px < bounds.Min.X || px >= bounds.Max.X {
			continue
		}
		for y := yStart; y <= yEnd; y++ {
			blendPixel(img, px, y, c)
		}
	}
}

func drawHorizontalLine(img *image.RGBA, y, xStart, xEnd, thickness int, c color.RGBA) {
	half := thickness / 2
	bounds := img.Bounds()
	for t := -half; t <= half; t++ {
		py := y + t
		if py < bounds.Min.Y || py >= bounds.Max.Y {
			continue
		}
		for x := xStart; x <= xEnd; x++ {
			blendPixel(img, x, py, c)
		}
	}
}

func blendPixel(img *image.RGBA, x, y int, c color.RGBA) {
	if c.A == 255 {
		img.SetRGBA(x, y, c)
		return
	}
	bg := img.RGBAAt(x, y)
	a := float64(c.A) / 255.0
	inv := 1.0 - a
	img.SetRGBA(x, y, color.RGBA{
		R: uint8(float64(c.R)*a + float64(bg.R)*inv),
		G: uint8(float64(c.G)*a + float64(bg.G)*inv),
		B: uint8(float64(c.B)*a + float64(bg.B)*inv),
		A: 255,
	})
}

func savePNG(img *image.RGBA, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de saída: %w", err)
	}
	defer f.Close()
	return png.Encode(f, img)
}
