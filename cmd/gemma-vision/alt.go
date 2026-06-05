package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/Siddhant-K-code/gemma-vision/internal/ollama"
)

var (
	altOutputJSON bool
	altPatchMDX   bool
	altOutputFile string
)

var altCmd = &cobra.Command{
	Use:   "alt <image-or-dir> [image-or-dir...]",
	Short: "Generate alt text for images",
	Long: `Walk one or more image files or directories and generate descriptive
alt text for each using gemma4:12b vision.

Output modes:
  default   Print image path + alt text to stdout
  --json    Write a JSON map of { path: altText } to stdout or --out file
  --patch   Rewrite MDX/HTML files in-place, filling empty alt="" attributes`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAlt,
}

func init() {
	altCmd.Flags().BoolVar(&altOutputJSON, "json", false, "Output results as JSON")
	altCmd.Flags().BoolVar(&altPatchMDX, "patch", false, "Patch alt=\"\" in MDX/HTML files in-place")
	altCmd.Flags().StringVar(&altOutputFile, "out", "", "Write JSON output to file instead of stdout")
}

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".gif": true, ".webp": true, ".bmp": true,
}

func collectImages(paths []string) ([]string, error) {
	var images []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			err = filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				if imageExts[strings.ToLower(filepath.Ext(path))] {
					images = append(images, path)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else if imageExts[strings.ToLower(filepath.Ext(p))] {
			images = append(images, p)
		}
	}
	return images, nil
}

const altPrompt = `Describe this image in one concise sentence suitable for use as HTML alt text.
Focus on the visual content and its meaning. Do not start with "This image shows" or "An image of".
Output only the alt text, nothing else.`

func runAlt(cmd *cobra.Command, args []string) error {
	client := ollama.New(ollamaHost, ollamaModel)

	images, err := collectImages(args)
	if err != nil {
		return err
	}
	if len(images) == 0 {
		return fmt.Errorf("no images found in the provided paths")
	}

	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	bar := progressbar.NewOptions(len(images),
		progressbar.OptionSetDescription("Generating alt text"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer: "█", SaucerPadding: "░", BarStart: "│", BarEnd: "│",
		}),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
	)

	results := make(map[string]string, len(images))

	for _, img := range images {
		// Spinner while model is thinking
		spinner := progressbar.NewOptions(-1,
			progressbar.OptionSetDescription(fmt.Sprintf("  Analysing %s", filepath.Base(img))),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer: "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏", SaucerPadding: " ", BarStart: "", BarEnd: "",
			}),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionSetWriter(os.Stderr),
		)

		done := make(chan struct{})
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					_ = spinner.Add(1)
					time.Sleep(80 * time.Millisecond)
				}
			}
		}()

		start := time.Now()
		alt, err := client.Ask(altPrompt, []string{img})
		close(done)
		_ = spinner.Clear()

		elapsed := time.Since(start).Round(time.Millisecond)

		if err != nil {
			red.Fprintf(os.Stderr, "✗ %s: %v\n", filepath.Base(img), err)
			_ = bar.Add(1)
			continue
		}
		alt = strings.TrimSpace(alt)
		results[img] = alt
		_ = bar.Add(1)

		if !altOutputJSON && !altPatchMDX {
			bold.Printf("\n%s\n", img)
			green.Printf("  %s\n", alt)
			color.New(color.Faint).Fprintf(os.Stderr, "  (%s)\n", elapsed)
		}
	}

	if altOutputJSON {
		data, _ := json.MarshalIndent(results, "", "  ")
		if altOutputFile != "" {
			if err := os.WriteFile(altOutputFile, data, 0644); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Wrote %s\n", altOutputFile)
		} else {
			fmt.Println(string(data))
		}
	}

	if altPatchMDX {
		if err := patchMDXFiles(args, results); err != nil {
			return err
		}
	}

	return nil
}

// patchMDXFiles walks MDX/HTML files and fills empty alt="" with generated alt text.
// Matches both <img src="path" alt=""> and ![](path) markdown image syntax.
var (
	htmlImgRe = regexp.MustCompile(`(<img[^>]+src=["'])([^"']+)(["'][^>]*alt=["'])["']`)
	mdImgRe   = regexp.MustCompile(`!\[]\(([^)]+)\)`)
)

func patchMDXFiles(roots []string, alts map[string]string) error {
	green := color.New(color.FgGreen)
	mdxExts := map[string]bool{".mdx": true, ".md": true, ".html": true, ".svelte": true}

	// Build a basename → alt map for fuzzy matching
	baseAlts := make(map[string]string, len(alts))
	for path, alt := range alts {
		baseAlts[filepath.Base(path)] = alt
	}

	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			if !mdxExts[strings.ToLower(filepath.Ext(path))] {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			original := string(content)
			patched := original

			// Patch HTML <img alt="">
			patched = htmlImgRe.ReplaceAllStringFunc(patched, func(m string) string {
				sub := htmlImgRe.FindStringSubmatch(m)
				if len(sub) < 4 {
					return m
				}
				src := sub[2]
				if alt, ok := baseAlts[filepath.Base(src)]; ok {
					return sub[1] + src + sub[3] + alt + `"`
				}
				return m
			})

			// Patch Markdown ![]()
			patched = mdImgRe.ReplaceAllStringFunc(patched, func(m string) string {
				sub := mdImgRe.FindStringSubmatch(m)
				if len(sub) < 2 {
					return m
				}
				src := sub[1]
				if alt, ok := baseAlts[filepath.Base(src)]; ok {
					return fmt.Sprintf("![%s](%s)", alt, src)
				}
				return m
			})

			if patched != original {
				if err := os.WriteFile(path, []byte(patched), 0644); err != nil {
					return err
				}
				green.Printf("Patched %s\n", path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
