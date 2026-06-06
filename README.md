# gemma-vision

Local image intelligence CLI powered by [gemma4](https://ollama.com/library/gemma4) via Ollama.

```bash
gemma-vision alt ./screenshots/          # generate alt text for every image
gemma-vision outline whiteboard.jpg      # turn a photo into a blog post outline
```

No API keys. No cloud. Runs entirely on your machine.

![Video demo](/gemma-vision-video-demo.mp4)


## Requirements

- [Ollama](https://ollama.com) running locally
- A vision model pulled — `ollama pull gemma4:12b` (or `gemma4:e4b` for 8GB RAM)

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/Siddhant-K-code/gemma-vision/main/install.sh | bash
```

Installs to `/usr/local/bin`. See [other install options](#other-install-options) below.

## Usage

### `alt` — Generate alt text for images

```bash
gemma-vision alt hero.png                        # single image
gemma-vision alt ./screenshots/                  # entire directory
gemma-vision alt ./images/ --json                # output as JSON
gemma-vision alt ./images/ --json --out alts.json
gemma-vision alt ./src/content/ --patch          # patch MDX/HTML files in-place
```

<details>
<summary>What <code>--patch</code> does</summary>

Walks MDX, Markdown, HTML, and Svelte files and fills in missing alt text:

- `![](image.png)` → `![Generated alt text](image.png)`
- `<img src="image.png" alt="">` → `<img src="image.png" alt="Generated alt text">`

</details>

### `outline` — Generate a blog post outline from images

```bash
gemma-vision outline whiteboard.jpg              # single photo
gemma-vision outline s1.png s2.png s3.png        # slide sequence
gemma-vision outline diagram.png --out outline.md
gemma-vision outline diagram.png \
  --audience "startup founders" \
  --tone "direct and opinionated"
```

<details>
<summary>Output format</summary>

```markdown
## Title Options
## Hook
## Sections
### Section Title
- key point
## Conclusion
## CTA
## SEO Keywords
```

</details>

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--model` | `gemma4:12b` | Ollama model (`GEMMA_MODEL` env var) |
| `--host` | `http://localhost:11434` | Ollama host (`OLLAMA_HOST` env var) |

<details>
<summary>All flags</summary>

**`alt`**

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |
| `--out <file>` | Write JSON to file |
| `--patch` | Patch MDX/HTML files in-place |

**`outline`**

| Flag | Default | Description |
|------|---------|-------------|
| `--audience` | `software engineers` | Target audience |
| `--tone` | `technical but conversational` | Writing tone |
| `--out <file>` | — | Write to markdown file |

</details>

## Other install options

<details>
<summary>Via Go / build from source / manual download</summary>

**Via Go:**
```bash
go install github.com/Siddhant-K-code/gemma-vision/cmd/gemma-vision@latest
```

**Build from source:**
```bash
git clone https://github.com/Siddhant-K-code/gemma-vision
cd gemma-vision
go build -o gemma-vision ./cmd/gemma-vision/
```

**Manual download:**
Grab the binary for your platform from the [Releases page](https://github.com/Siddhant-K-code/gemma-vision/releases), extract, and move to your `$PATH`.

**Custom install directory:**
```bash
INSTALL_DIR=~/.local/bin curl -fsSL https://raw.githubusercontent.com/Siddhant-K-code/gemma-vision/main/install.sh | bash
```

</details>
