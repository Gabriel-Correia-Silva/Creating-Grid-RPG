package colors

import (
	"bufio"
	"fmt"

	"grid-generator/internal/config"
	"grid-generator/internal/input"
)

var linePresets = map[int]config.RGB{
	1: {R: 200, G: 200, B: 200},
	2: {R: 0, G: 0, B: 0},
	3: {R: 255, G: 255, B: 255},
}

var backgroundPresets = map[int]config.RGB{
	2: {R: 30, G: 30, B: 30},
}

func LineFromPreset(choice int) config.RGB {
	if c, ok := linePresets[choice]; ok {
		return c
	}
	return linePresets[1]
}

func BackgroundFromPreset(choice int) *config.RGB {
	if c, ok := backgroundPresets[choice]; ok {
		cp := c
		return &cp
	}
	return nil
}

func ReadRGB(r *bufio.Reader, label string) config.RGB {
	fmt.Printf("Informe a cor %s em RGB (0-255):\n", label)
	fmt.Print("  R (vermelho): ")
	rv := input.ReadInt(r)
	fmt.Print("  G (verde):    ")
	gv := input.ReadInt(r)
	fmt.Print("  B (azul):     ")
	bv := input.ReadInt(r)
	return config.RGB{
		R: clamp(rv, 0, 255),
		G: clamp(gv, 0, 255),
		B: clamp(bv, 0, 255),
	}
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
