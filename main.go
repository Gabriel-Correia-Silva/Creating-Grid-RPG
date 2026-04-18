package main

import (
	"bufio"
	"fmt"
	"os"

	"grid-generator/internal/cli"
	"grid-generator/internal/config"
	"grid-generator/internal/pdf"
	"grid-generator/internal/raster"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	cli.PrintHeader()

	cfg, err := cli.GatherConfig(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERRO] %v\n", err)
		os.Exit(1)
	}

	cli.PrintSummary(cfg)

	if err := dispatch(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "[ERRO] %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[SUCESSO] Arquivo '%s' gerado com sucesso!\n", cfg.OutputFile)
}

func dispatch(cfg config.GridConfig) error {
	if cfg.OutputPDF {
		fmt.Println("[INFO] Gerando PDF vetorial...")
		return pdf.Generate(cfg)
	}
	fmt.Println("[INFO] Gerando imagem PNG...")
	return raster.Generate(cfg)
}
