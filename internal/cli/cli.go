package cli

import (
	"bufio"
	"fmt"
	"math"
	"strings"

	"grid-generator/internal/colors"
	"grid-generator/internal/config"
	"grid-generator/internal/images"
	"grid-generator/internal/input"
)

func GatherConfig(r *bufio.Reader) (config.GridConfig, error) {
	cfg := config.GridConfig{LineOpacity: 1.0}

	cfg.OutputPDF = promptOutputFormat(r)

	if err := promptBackgroundImage(r, &cfg); err != nil {
		return config.GridConfig{}, err
	}

	if cfg.BackgroundImage == "" {
		promptPaperSize(r, &cfg)
	}

	cfg.GridCellCm = promptGridCellSize(r)
	cfg.LineWidth = promptLineWidth(r)
	cfg.LineColor = promptLineColor(r)
	cfg.LineOpacity = promptLineOpacity(r)

	if cfg.BackgroundImage == "" {
		cfg.BackgroundColor = promptBackgroundColor(r)
	}

	cfg.MarginCm = promptMargin(r)
	cfg.OutputFile = promptOutputFile(r, cfg.OutputPDF)

	return cfg, nil
}

func promptOutputFormat(r *bufio.Reader) bool {
	fmt.Println("--- 1. FORMATO DE SAÍDA ---")
	fmt.Println("1: PDF (vetorial) | 2: Imagem PNG (raster)")
	fmt.Print("Escolha: ")
	return input.ReadInt(r) != 2
}

func promptBackgroundImage(r *bufio.Reader, cfg *config.GridConfig) error {
	fmt.Println("\n--- 2. IMAGEM DE FUNDO ---")
	fmt.Println("1: Não | 2: Sim (informar caminho)")
	fmt.Print("Escolha: ")
	if input.ReadInt(r) != 2 {
		return nil
	}

	fmt.Print("Caminho da imagem: ")
	path := input.ReadLine(r)

	w, h, err := images.Size(path)
	if err != nil {
		return fmt.Errorf("não foi possível ler a imagem '%s': %w", path, err)
	}

	cfg.BackgroundImage = path
	cfg.PageWidth = float64(w)
	cfg.PageHeight = float64(h)

	fmt.Printf("[INFO] Imagem carregada: %dx%d px\n", w, h)
	fmt.Println("[INFO] Grid gerado na proporção da imagem.")
	return nil
}

func promptPaperSize(r *bufio.Reader, cfg *config.GridConfig) {
	fmt.Println("\n--- 3. PAPEL ---")
	fmt.Println("1: A4 | 2: A3 | 3: A1 | 4: Customizado (cm)")
	fmt.Print("Escolha: ")
	choice := input.ReadInt(r)

	if size, ok := config.PaperPresets[choice]; ok {
		cfg.PageWidth = size[0]
		cfg.PageHeight = size[1]
	} else if choice == 4 {
		fmt.Print("LARGURA (cm): ")
		cfg.PageWidth = input.ReadFloat(r) * config.CmToPoints
		fmt.Print("ALTURA (cm): ")
		cfg.PageHeight = input.ReadFloat(r) * config.CmToPoints
	} else {
		cfg.PageWidth = config.PaperPresets[1][0]
		cfg.PageHeight = config.PaperPresets[1][1]
	}

	fmt.Println("\nOrientação -> 1: Retrato | 2: Paisagem")
	fmt.Print("Escolha: ")
	if input.ReadInt(r) == 2 {
		cfg.PageWidth, cfg.PageHeight = cfg.PageHeight, cfg.PageWidth
	}
}

func promptGridCellSize(r *bufio.Reader) float64 {
	fmt.Println("\n--- 4. DIMENSÕES DO GRID ---")
	fmt.Print("Tamanho do quadrado em cm (ex: 1.5): ")
	return input.ReadFloat(r)
}

func promptLineWidth(r *bufio.Reader) float64 {
	fmt.Println("\n--- 5. ESPESSURA DA LINHA ---")
	fmt.Print("Largura em pontos (0.3=fina, 1.0=normal, 2.0=grossa): ")
	return input.ReadFloat(r)
}

func promptLineColor(r *bufio.Reader) config.RGB {
	fmt.Println("\n--- 6. COR DA LINHA ---")
	fmt.Println("1: Cinza Claro | 2: Preto | 3: Azul | 4: Vermelho")
	fmt.Println("5: Verde | 6: Branco | 7: Customizada (RGB)")
	fmt.Print("Escolha: ")
	choice := input.ReadInt(r)
	if choice == 7 {
		return colors.ReadRGB(r, "linha")
	}
	return colors.LineFromPreset(choice)
}

func promptLineOpacity(r *bufio.Reader) float64 {
	fmt.Println("\n--- 7. OPACIDADE DA LINHA ---")
	fmt.Print("Opacidade (0.0=invisível, 1.0=opaca) [padrão: 1.0]: ")
	s := input.ReadLine(r)
	if s == "" {
		return 1.0
	}
	return math.Max(0, math.Min(1, input.ParseFloat(s)))
}

func promptBackgroundColor(r *bufio.Reader) *config.RGB {
	fmt.Println("\n--- 8. COR DE FUNDO ---")
	fmt.Println("1: Branco (padrão) | 2: Amarelo Claro | 3: Azul Claro")
	fmt.Println("4: Preto | 5: Customizada (RGB)")
	fmt.Print("Escolha: ")
	choice := input.ReadInt(r)
	if choice == 5 {
		c := colors.ReadRGB(r, "fundo")
		return &c
	}
	return colors.BackgroundFromPreset(choice)
}

func promptMargin(r *bufio.Reader) float64 {
	fmt.Println("\n--- 9. MARGEM ---")
	fmt.Print("Margem em cm (0 = sem margem) [padrão: 0]: ")
	s := input.ReadLine(r)
	if s == "" {
		return 0
	}
	return input.ParseFloat(s)
}

func promptOutputFile(r *bufio.Reader, isPDF bool) string {
	ext := ".png"
	if isPDF {
		ext = ".pdf"
	}
	fmt.Printf("\nNome do arquivo (ex: meu_grid): ")
	name := input.ReadLine(r)
	for _, e := range []string{".pdf", ".png", ".jpg", ".bmp", ".gif", ".tiff", ".tif", ".webp"} {
		if strings.HasSuffix(name, e) {
			name = name[:len(name)-len(e)]
			break
		}
	}
	return name + ext
}

func PrintHeader() {
	fmt.Println("============================================")
	fmt.Println("   GERADOR DE GRID - MODO AVANÇADO v3.0    ")
	fmt.Println("   (Go Edition — Cross-platform)           ")
	fmt.Println("============================================")
	fmt.Println()
}

func PrintSummary(cfg config.GridConfig) {
	fmt.Println("\n========= RESUMO =========")
	if cfg.OutputPDF {
		fmt.Println("Formato:       PDF (vetorial)")
	} else {
		fmt.Println("Formato:       PNG (imagem)")
	}
	fmt.Printf("Papel:         %.0f x %.0f pts\n", cfg.PageWidth, cfg.PageHeight)
	fmt.Printf("Quadrado:      %.2f cm\n", cfg.GridCellCm)
	fmt.Printf("Linha:         %.2f pts, RGB(%d,%d,%d), opacidade=%.2f\n",
		cfg.LineWidth, cfg.LineColor.R, cfg.LineColor.G, cfg.LineColor.B, cfg.LineOpacity)
	switch {
	case cfg.BackgroundImage != "":
		fmt.Printf("Fundo:         Imagem (%s)\n", cfg.BackgroundImage)
	case cfg.BackgroundColor != nil:
		fmt.Printf("Fundo:         RGB(%d,%d,%d)\n",
			cfg.BackgroundColor.R, cfg.BackgroundColor.G, cfg.BackgroundColor.B)
	default:
		fmt.Println("Fundo:         Branco (padrão)")
	}
	fmt.Printf("Margem:        %.2f cm\n", cfg.MarginCm)
	fmt.Printf("Arquivo:       %s\n", cfg.OutputFile)
	fmt.Println("==========================")
	fmt.Println()
}
