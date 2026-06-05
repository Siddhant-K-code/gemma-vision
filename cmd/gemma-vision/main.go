package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	ollamaHost  string
	ollamaModel string
)

var rootCmd = &cobra.Command{
	Use:   "gemma-vision",
	Short: "Local image intelligence CLI powered by gemma4:12b via Ollama",
	Long: color.New(color.Bold).Sprint("gemma-vision") + ` — run vision tasks locally using gemma4:12b.

Subcommands:
  alt      Generate alt text for images (single file, directory, or MDX patch)
  outline  Generate a blog post outline from one or more images

Examples:
  gemma-vision alt ./screenshots/
  gemma-vision alt ./img/hero.png --json
  gemma-vision alt ./src/content/ --patch
  gemma-vision outline ./whiteboard.jpg
  gemma-vision outline slide1.png slide2.png slide3.png --out outline.md`,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&ollamaHost, "host", envOr("OLLAMA_HOST", "http://localhost:11434"), "Ollama host URL")
	rootCmd.PersistentFlags().StringVar(&ollamaModel, "model", envOr("GEMMA_MODEL", "gemma4:12b"), "Ollama model to use")

	rootCmd.AddCommand(altCmd)
	rootCmd.AddCommand(outlineCmd)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error: ")+err.Error())
		os.Exit(1)
	}
}
