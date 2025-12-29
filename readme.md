# Font Lengths Generator

`fontlengths` is a Go tool that generates visual representations of character widths for all fonts installed on your system. It creates a PNG image for each font, displaying lines of repeated characters sorted by their pixel length.

This tool was inspired by a tweet about font character width visualization: https://twitter.com/FreyaHolmer/status/1776007184117063807

## Features

- **Automatic Font Discovery**: Scans your system for installed fonts using `go-findfont`.
- **Format Support**: Supports both TrueType (`.ttf`) and OpenType (`.otf`) fonts.
- **Visualization**:
  - Generates a line for each letter from 'A' to 'Z'.
  - Each line consists of 20 repetitions of the character (e.g., "AAAA...").
  - Lines are sorted by their total pixel length (width), allowing you to easily see which characters are widest and narrowest in a specific font.
  - Images are rendered at high resolution (16pt font size at 300 DPI).
- **Batch Processing**: Processes all found fonts in one go.

## Installation

### Prerequisites

- [Go](https://go.dev/dl/) (version 1.21 or later recommended)
- System fonts (standard on most operating systems)

### Build from Source

1.  Clone the repository:
    ```bash
    git clone https://github.com/arran4/fontlengths.git
    cd fontlengths
    ```

2.  Download dependencies:
    ```bash
    go mod download
    ```

3.  Build the binary:
    ```bash
    go build -o fontlengths ./cmd/generate
    ```

## Usage

Run the built executable:

```bash
./fontlengths
```

The tool will log its progress as it finds and processes fonts.

### Output

Generated images are saved in the `out/` directory within the project folder. The directory will be created automatically if it doesn't exist.

Each file is named after the font's full name, for example:
- `out/Arial.png`
- `out/Ubuntu Regular.png`
- `out/DejaVu Sans Mono.png`

**Example Output:**

![Ubuntu Regular.png](Ubuntu%20Regular.png)

*(Note: The above image is an example. Run the tool to generate images for your own fonts.)*

## Configuration

Currently, the tool operates with default settings hardcoded in `cmd/generate/main.go`. You can modify the source code to change:

- **Font Size**: Change the `16` in `GetTTFontFace` and `GetOTFontFace` calls.
- **DPI**: Change the `300` in the same calls.
- **Text Content**: Modify the loop in `CreateImage` (currently 'A'-'Z', repeated 20 times).
- **Colors**: Change `image.White` (background) or `image.Black` (text) in `CreateImage`.

## Troubleshooting

- **No fonts found**: Ensure your system has fonts installed in standard locations. The `go-findfont` library searches common font directories.
- **Permission errors**: If the tool crashes or fails to write files, check that you have write permissions for the `out/` directory.
- **Skipped fonts**: Some fonts might be corrupted or in an unsupported format. The tool logs these errors and skips to the next font.

## License

This project is open source. See the LICENSE file for details.

## Credits

- **Inspiration**: [Freya Holmér](https://twitter.com/FreyaHolmer)
- **Dependencies**:
  - `github.com/arran4/golang-wordwrap`
  - `github.com/flopp/go-findfont`
  - `golang.org/x/image`
  - `github.com/golang/freetype`
