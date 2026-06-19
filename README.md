# Go Base30 to Image

A professional CLI tool to convert jSignature signature data from Base30 format to PNG or JPG images. Built with Go and Cobra CLI framework.

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey.svg)](https://github.com/swantara-tech/go-base30-to-image)

## Features

- ✅ Decode jSignature Base30 format with delta encoding
- ✅ Convert single signatures from file or command line string
- ✅ Batch processing from CSV files
- ✅ Export to PNG (lossless) or JPG (configurable quality)
- ✅ Auto-sizing canvas with configurable dimensions
- ✅ White background with black strokes
- ✅ Configurable stroke width and anti-aliasing
- ✅ Verbose logging for debugging
- ✅ Comprehensive error handling
- ✅ Unit tests with coverage reporting
- ✅ Cross-platform support (Windows, Linux, macOS)
- ✅ Multi-platform binary build

## Installation

### From Source

```bash
git clone https://github.com/swantara-tech/go-base30-to-image.git
cd go-base30-to-image
go mod tidy
go build -o jsign-convert
```

### Using Make

```bash
make deps    # Install dependencies
make build   # Build binary
```

## Download Pre-built Binaries

Download from [GitHub Releases](https://github.com/swantara-tech/go-base30-to-image/releases):
- **Windows:** `jsign-convert-v1.0.0-windows-amd64.exe`
- **Linux:** `jsign-convert-v1.0.0-linux-amd64`
- **macOS:** `jsign-convert-v1.0.0-darwin-amd64`

### Verify Checksums

Each release includes `checksums.txt` with SHA256 hashes. Verify your download:

**Linux/macOS:**
```bash
sha256sum -c checksums.txt
```

**Windows (PowerShell):**
```powershell
Get-FileHash jsign-convert-v1.0.0-windows-amd64.exe -Algorithm SHA256
```

## Usage

### Single Conversion

#### From File

**Linux/macOS:**
```bash
./jsign-convert convert \
  --input signature.txt \
  --output signature.png
```

**Windows (PowerShell):**
```powershell
.\jsign-convert.exe convert --input signature.txt --output signature.png
```

#### From String

**Linux/macOS:**
```bash
./jsign-convert convert \
  --base30 "5K247669cffhlo1vmhc9852" \
  --output signature.png
```

**Windows (PowerShell):**
```powershell
.\jsign-convert.exe convert --base30 "5K247669cffhlo1vmhc9852" --output signature.png
```

#### With Custom Dimensions

**Linux/macOS:**
```bash
./jsign-convert convert \
  --input signature.txt \
  --output signature.png \
  --width 800 \
  --height 300
```

**Windows (PowerShell):**
```powershell
.\jsign-convert.exe convert --input signature.txt --output signature.png --width 800 --height 300
```

#### JPG Output

**Linux/macOS:**
```bash
./jsign-convert convert \
  --input signature.txt \
  --output signature.jpg \
  --format jpg \
  --quality 95
```

**Windows (PowerShell):**
```powershell
.\jsign-convert.exe convert --input signature.txt --output signature.jpg --format jpg --quality 95
```

### Batch Conversion

**Linux/macOS:**
```bash
./jsign-convert batch \
  --csv signatures.csv \
  --output-dir ./results
```

**Windows (PowerShell):**
```powershell
.\jsign-convert.exe batch --csv signatures.csv --output-dir .\results
```

#### Batch with JPG Output

**Linux/macOS:**
```bash
./jsign-convert batch \
  --csv signatures.csv \
  --output-dir ./results \
  --format jpg \
  --quality 90
```

**Windows (PowerShell):**
```powershell
.\jsign-convert.exe batch --csv signatures.csv --output-dir .\results --format jpg --quality 90
```

### Verbose Mode

**Linux/macOS:**
```bash
./jsign-convert convert --input signature.txt --verbose
```

**Windows (PowerShell):**
```powershell
.\jsign-convert.exe convert --input signature.txt --verbose
```

## Input Format

### Supported Formats

The tool accepts both formats:

1. **With prefix:**
   ```
   image/jsignature;base30,5K247669cffhlo1vmhc9852Z346...
   ```

2. **Without prefix:**
   ```
   5K247669cffhlo1vmhc9852Z346...
   ```

### CSV Format for Batch Processing

```csv
id,signature
1,image/jsignature;base30,5K247669cffhlo1vmhc9852...
2,image/jsignature;base30,7L358770dglimp2wnid0963...
```

## Command Reference

### Convert Command

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--input` | `-i` | Input file containing signature data | - |
| `--output` | `-o` | Output file path | `signature.png` |
| `--base30` | - | Base30 signature string | - |
| `--format` | - | Output format (png or jpg) | `png` |
| `--quality` | - | JPG quality (1-100) | `95` |
| `--width` | - | Output width (0 for auto) | `0` |
| `--height` | - | Output height (0 for auto) | `0` |
| `--verbose` | `-v` | Enable verbose logging | `false` |

### Batch Command

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--csv` | `-c` | Input CSV file | - |
| `--output-dir` | `-d` | Output directory | `./results` |
| `--format` | - | Output format (png or jpg) | `png` |
| `--quality` | - | JPG quality (1-100) | `95` |
| `--verbose` | `-v` | Enable verbose logging | `false` |

## Project Structure

```
cmd/
├── root.go          # Root command and CLI setup
├── convert.go       # Single conversion command
├── batch.go         # Batch conversion command

internal/
├── decoder/
│   ├── base30.go    # Base30 encoding/decoding
│   ├── parser.go    # Signature stroke parser
│
├── renderer/
│   ├── image.go     # Image rendering engine
│   ├── jpg.go       # JPG encoding
│
├── models/
│   └── signature.go # Data models

pkg/
└── utils/
    └── file.go      # File and CSV utilities

main.go              # Entry point
```

## Development

### Build

```bash
make build              # Build for current platform
make build-all          # Build for Windows, Linux, and macOS
```

### Run Tests

```bash
make test               # Run tests with coverage
make test-coverage      # Generate HTML coverage report
make test-coverage-check # Check coverage percentage
```

### Code Quality

```bash
make fmt                # Format code with gofmt
make lint               # Run golangci-lint
```

### Utilities

```bash
make examples           # Create example signature files
make clean              # Remove build artifacts
make deps               # Install/update dependencies
make run ARGS="..."     # Run without building
```

### Build Multi-Platform Binaries

```bash
make build-all
```

Output:
```
build/
├── jsign-convert-linux-amd64
├── jsign-convert-windows-amd64.exe
└── jsign-convert-darwin-amd64
```

## Testing

Run all tests:

```bash
go test -v ./...
```

Run tests with coverage:

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Error Handling

The tool handles various error conditions:

- ❌ Invalid Base30 data
- ❌ Empty signature data
- ❌ Corrupted signature format
- ❌ Invalid output paths
- ❌ Missing required flags

Returns non-zero exit code on failure.

## Examples

### Example 1: Simple Conversion

**Linux/macOS:**
```bash
# Create a signature file
echo "image/jsignature;base30,5K247669cffhlo1vmhc9852" > signature.txt

# Convert to PNG
./jsign-convert convert --input signature.txt --output signature.png
```

**Windows (PowerShell):**
```powershell
# Create a signature file
echo "image/jsignature;base30,5K247669cffhlo1vmhc9852" | Out-File -Encoding ASCII signature.txt

# Convert to PNG
.\jsign-convert.exe convert --input signature.txt --output signature.png
```

### Example 2: Batch Processing

**Linux/macOS:**
```bash
# Create CSV file
cat > signatures.csv << EOF
id,signature
1,image/jsignature;base30,5K247669cffhlo1vmhc9852
2,image/jsignature;base30,7L358770dglimp2wnid0963
EOF

# Batch convert
./jsign-convert batch --csv signatures.csv --output-dir ./results
```

**Windows (PowerShell):**
```powershell
# Create CSV file
@"
id,signature
1,image/jsignature;base30,5K247669cffhlo1vmhc9852
2,image/jsignature;base30,7L358770dglimp2wnid0963
"@ | Out-File -Encoding ASCII signatures.csv

# Batch convert
.\jsign-convert.exe batch --csv signatures.csv --output-dir .\results
```

Output:
```
results/
├── 1.png
└── 2.png
```

## Making Releases

### Create a New Release

```bash
# Create and push a new tag
git tag -a v1.0.0 -m "Initial release"
git push origin v1.0.0
```

GitHub Actions will automatically:
- ✅ Run tests
- ✅ Build binaries for all platforms (Windows, Linux, macOS)
- ✅ Generate SHA256 checksums
- ✅ Create GitHub Release with all binaries
- ✅ Generate changelog

### Versioning

This project follows [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality (backward compatible)
- **PATCH** version for backward compatible bug fixes

Example: `v1.0.0`, `v1.1.0`, `v1.1.1`, `v2.0.0`

## CI/CD

This project uses GitHub Actions for continuous integration:

- ✅ Automated testing on push and PR
- ✅ Multi-platform builds (Linux, Windows, macOS)
- ✅ Code coverage reporting
- ✅ Linting with golangci-lint

See [.github/workflows/ci.yml](.github/workflows/ci.yml) for configuration.

## Requirements

- **Go:** 1.24 or higher
- **Build Tools:** GNU Make (optional but recommended)
- **Runtime:** None (self-contained binary)
- **Dependencies:** No external runtime dependencies

## Platform Support

| Platform | Architecture | Status |
|----------|--------------|--------|
| Linux    | amd64        | ✅     |
| Windows  | amd64        | ✅     |
| macOS    | amd64        | ✅     |

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make build` | Build binary for current platform |
| `make build-all` | Build for Windows, Linux, and macOS |
| `make test` | Run tests with coverage |
| `make test-coverage` | Generate HTML coverage report |
| `make test-coverage-check` | Check coverage percentage |
| `make fmt` | Format code |
| `make lint` | Run linter |
| `make clean` | Clean build artifacts |
| `make deps` | Install dependencies |
| `make run ARGS="..."` | Run application |
| `make examples` | Create example files |
| `make help` | Show all available commands |

## License

MIT License

## Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests: `make test`
4. Format code: `make fmt`
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Follow Go conventions (effective-go)
- Run `make lint` before submitting PR
- Update documentation as needed

## Support

- 📖 **Documentation:** [README.md](README.md)
- 🐛 **Bug Reports:** [GitHub Issues](https://github.com/swantara-tech/go-base30-to-image/issues)
- 💡 **Feature Requests:** [GitHub Issues](https://github.com/swantara-tech/go-base30-to-image/issues)
- 📧 **Contact:** [Swantara Tech](https://github.com/swantara-tech)

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

Made with ❤️ by [Swantara Tech](https://github.com/swantara-tech)
