package config

const CmToPoints = 72.0 / 2.54

var PaperPresets = map[int][2]float64{
	1: {595.28, 841.89},
	2: {841.89, 1190.55},
	3: {1683.78, 2383.94},
}

type RGB struct {
	R, G, B int
}

type GridConfig struct {
	OutputPDF       bool
	OutputFile      string
	PageWidth       float64
	PageHeight      float64
	GridCellCm      float64
	LineWidth       float64
	LineColor       RGB
	LineOpacity     float64
	BackgroundColor *RGB
	BackgroundImage string
	MarginCm        float64
}
