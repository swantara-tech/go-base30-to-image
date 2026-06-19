# JSignature Base30 Converter

## Project Overview

Buat aplikasi command line menggunakan Golang untuk mengkonversi data tanda tangan jSignature format `image/jsignature;base30` menjadi file PNG atau JPG.

Aplikasi harus dapat digunakan untuk memproses satu data maupun batch processing dari file CSV atau database.

---

# Objectives

1. Decode format `image/jsignature;base30`
2. Mengubah data stroke menjadi koordinat X dan Y
3. Render stroke ke canvas image
4. Simpan hasil sebagai PNG
5. Opsional convert PNG menjadi JPG
6. Mendukung batch processing

---

# Technical Requirements

## Language

* Golang 1.24+
* Tanpa dependency frontend
* CLI Application

## Supported Input

Format:

```text
image/jsignature;base30,5K247669cffhlo1vmhc9852Z346aaddfeba50Y247cefgfdc72...
```

atau

```text
5K247669cffhlo1vmhc9852Z346aaddfeba50Y247cefgfdc72...
```

---

# Features

## Single Convert

Command:

```bash
jsign-convert \
--input signature.txt \
--output signature.png
```

Output:

```text
signature.png
```

---

## Direct String Convert

Command:

```bash
jsign-convert \
--base30 "5K247669cffhlo1vmhc9852..." \
--output signature.png
```

---

## Batch Convert

Command:

```bash
jsign-convert batch \
--csv signatures.csv \
--output-dir ./results
```

CSV Format:

```csv
id,signature
1,image/jsignature;base30,...
2,image/jsignature;base30,...
```

Output:

```text
results/
├── 1.png
├── 2.png
├── 3.png
```

---

# Core Decoder Requirements

Implementasikan parser jSignature Base30.

Harus mendukung:

* Multiple stroke
* Delta encoding
* Base30 character decoding
* Coordinate reconstruction
* Pen up / pen down handling

Referensi format:

https://github.com/brinley/jSignature

Jangan menggunakan browser.

Implementasi harus native Golang.

---

# Rendering Engine

Gunakan package:

```go
image
image/color
image/draw
```

atau

```go
github.com/fogleman/gg
```

Requirements:

* Background putih
* Stroke hitam
* Auto sizing canvas
* Margin 20px
* Anti alias jika memungkinkan

---

# PNG Export

Default output:

```bash
signature.png
```

Resolution:

```text
Auto Fit
```

atau

```bash
--width 800
--height 300
```

---

# JPG Export

Command:

```bash
--format jpg
```

Quality:

```bash
--quality 95
```

---

# Project Structure

```text
cmd/
├── root.go
├── convert.go
├── batch.go

internal/
├── decoder/
│   ├── base30.go
│   ├── parser.go
│
├── renderer/
│   ├── png.go
│   ├── jpg.go
│
├── models/
│   └── signature.go

pkg/
└── utils/

main.go
```

---

# Error Handling

Handle:

* Invalid base30 data
* Empty signature
* Corrupted signature
* Invalid output path

Return exit code non-zero jika gagal.

---

# Logging

Gunakan:

```go
log/slog
```

Support:

```bash
--verbose
```

---

# Unit Test

Coverage minimal:

* Base30 decoder
* Stroke parser
* Renderer
* PNG export

Coverage target:

```text
>= 80%
```

---

# Deliverables

Generate:

1. Full Golang source code
2. Unit tests
3. README.md
4. Example input
5. Example output PNG
6. Makefile
7. GitHub Actions CI

Project harus bisa dijalankan dengan:

```bash
go mod tidy
go build
```

Tanpa dependency browser atau JavaScript runtime.
