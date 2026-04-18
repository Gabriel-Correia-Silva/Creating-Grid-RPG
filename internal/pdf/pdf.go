package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"

	"grid-generator/internal/config"
	"grid-generator/internal/images"
)

const mmPerPoint = 25.4 / 72.0

func Generate(cfg config.GridConfig) error {
	widthMM := cfg.PageWidth / config.CmToPoints * 10.0
	heightMM := cfg.PageHeight / config.CmToPoints * 10.0

	doc := newDocument(widthMM, heightMM)

	if err := applyBackground(doc, cfg, widthMM, heightMM); err != nil {
		return err
	}

	applyLineStyle(doc, cfg)
	drawGrid(doc, cfg, widthMM, heightMM)

	return doc.OutputFileAndClose(cfg.OutputFile)
}

func newDocument(widthMM, heightMM float64) *gofpdf.Fpdf {
	doc := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size:    gofpdf.SizeType{Wd: widthMM, Ht: heightMM},
	})
	doc.SetMargins(0, 0, 0)
	doc.SetAutoPageBreak(false, 0)
	doc.AddPage()
	return doc
}

func applyBackground(doc *gofpdf.Fpdf, cfg config.GridConfig, widthMM, heightMM float64) error {
	if cfg.BackgroundImage != "" {
		return embedImage(doc, cfg.BackgroundImage, widthMM, heightMM)
	}
	if cfg.BackgroundColor != nil {
		doc.SetFillColor(cfg.BackgroundColor.R, cfg.BackgroundColor.G, cfg.BackgroundColor.B)
		doc.Rect(0, 0, widthMM, heightMM, "F")
	}
	return nil
}

func embedImage(doc *gofpdf.Fpdf, path string, widthMM, heightMM float64) error {
	imgPath, tempFile, err := resolveImagePath(path)
	if err != nil {
		return fmt.Errorf("erro ao preparar imagem de fundo: %w", err)
	}
	if tempFile != "" {
		defer os.Remove(tempFile)
	}

	opt := gofpdf.ImageOptions{ImageType: resolveImageType(imgPath), ReadDpi: true}
	doc.ImageOptions(imgPath, 0, 0, widthMM, heightMM, false, opt, 0, "")
	return nil
}

func resolveImagePath(path string) (imgPath, tempToDelete string, err error) {
	if images.IsPDFNative(path) {
		return path, "", nil
	}
	tmp, err := images.ConvertToPNGTemp(path)
	if err != nil {
		return "", "", err
	}
	return tmp, tmp, nil
}

func resolveImageType(path string) string {
	if strings.ToLower(filepath.Ext(path)) == ".png" {
		return "PNG"
	}
	return "JPG"
}

func applyLineStyle(doc *gofpdf.Fpdf, cfg config.GridConfig) {
	if cfg.LineOpacity < 1.0 {
		doc.SetAlpha(cfg.LineOpacity, "Normal")
	}
	doc.SetLineWidth(cfg.LineWidth * mmPerPoint)
	doc.SetDrawColor(cfg.LineColor.R, cfg.LineColor.G, cfg.LineColor.B)
}

func drawGrid(doc *gofpdf.Fpdf, cfg config.GridConfig, widthMM, heightMM float64) {
	spacingMM := cfg.GridCellCm * 10.0
	marginMM := cfg.MarginCm * 10.0

	for x := marginMM; x <= widthMM-marginMM; x += spacingMM {
		doc.Line(x, marginMM, x, heightMM-marginMM)
	}
	for y := marginMM; y <= heightMM-marginMM; y += spacingMM {
		doc.Line(marginMM, y, widthMM-marginMM, y)
	}
}
