# gemma-vision

Local image intelligence CLI powered by [gemma4:12b](https://ollama.com/library/gemma4) via Ollama.

Two subcommands:
- **`alt`** — generate alt text for images, output as JSON, or patch MDX/HTML files in-place
- **`outline`** — generate a structured blog post outline from one or more images

Everything runs locally. No API keys, no cloud.

## Requirements

- [Ollama](https://ollama.com) running locally
- `gemma4:12b` pulled: `ollama pull gemma4:12b`

## Install

```bash
go install github.com/Siddhant-K-code/gemma-vision/cmd/gemma-vision@latest
```

Or build from source:

```bash
git clone https://github.com/Siddhant-K-code/gemma-vision
cd gemma-vision
go build -o gemma-vision ./cmd/gemma-vision/
```

## Usage

### `alt` — Generate alt text for images

```bash
# Single image
gemma-vision alt hero.png

# Entire directory
gemma-vision alt ./screenshots/

# Output as JSON
gemma-vision alt ./images/ --json

# Write JSON to file
gemma-vision alt ./images/ --json --out alts.json

# Patch empty alt="" in MDX/HTML files in-place
gemma-vision alt ./src/content/ --patch
```

The `--patch` flag walks MDX, Markdown, HTML, and Svelte files and fills:
- `![](image.png)` → `![Generated alt text](image.png)`
- `<img src="image.png" alt="">` → `<img src="image.png" alt="Generated alt text">`

### `outline` — Generate a blog post outline from images

```bash
# Single whiteboard photo
gemma-vision outline whiteboard.jpg

# Multiple slides treated as a sequence
gemma-vision outline slide1.png slide2.png slide3.png

# Custom audience and tone
gemma-vision outline diagram.png --audience "startup founders" --tone "direct and opinionated"

# Write outline to a markdown file
gemma-vision outline whiteboard.jpg --out outline.md
```

Output format:

```markdown
## Title Options
1. ...
2. ...
3. ...

## Hook
...

## Sections
### Section Title
- key point
- key point

## Conclusion
...

## CTA
...

## SEO Keywords
...
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | `http://localhost:11434` | Ollama host (or set `OLLAMA_HOST`) |
| `--model` | `gemma4:12b` | Model to use (or set `GEMMA_MODEL`) |

### `alt` flags

| Flag | Description |
|------|-------------|
| `--json` | Output results as JSON |
| `--out <file>` | Write JSON to file |
| `--patch` | Patch MDX/HTML files in-place |

### `outline` flags

| Flag | Default | Description |
|------|---------|-------------|
| `--audience` | `software engineers` | Target audience |
| `--tone` | `technical but conversational` | Writing tone |
| `--out <file>` | — | Write outline to markdown file |

## Use a different model

```bash
gemma-vision --model gemma4:e4b alt ./images/
GEMMA_MODEL=llava:13b gemma-vision outline diagram.png
```
