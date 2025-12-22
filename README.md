# Ollama Readme Generator

This repository contains a small Go application that automatically generates a `README.md` file for a codebase.  
It reads every source file (respecting the repository’s `.gitignore`), appends a user‑supplied prompt, sends the combined text to an Ollama model, and streams the model’s response directly to stdout. The output can then be redirected to `README.md`.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Architecture](#architecture)
  - [Main package](#main-package)
  - [File processing](#file-processing)
  - [Ollama client](#ollama-client)
- [Functions](#functions)
  - [parseGitIgnore](#parsegitignore)
  - [isIgnored](#isignored)
  - [ReadData](#readdata)
  - [AskOllama](#askollama)
- [Logging & Debug](#logging--debug)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

The tool performs the following steps:

1. **Read a prompt** from a file (default: `Prompt.md`).  
2. **Scan the current repository** (starting from the working directory) for all files, **excluding** paths that match rules defined in `.gitignore`.  
3. **Concatenate** the prompt with the collected file contents.  
4. **Send** the combined payload to an Ollama server (`http://localhost:11434/api/generate`) via an HTTP POST request.  
5. **Stream** the JSON chunks returned by the Ollama API, printing the `response` field of each chunk to stdout.  
6. **Redirect** the stdout to a file to create a fully‑formed `README.md`.

The code is intentionally straightforward: no concurrency, no complex templating, and a clear separation of concerns between file handling and network communication.

---

## Features

| Feature | Description |
|---------|-------------|
| **Automatic file collection** | Recursively walks the repository, respecting `.gitignore`. |
| **Prompt injection** | User‑supplied text can shape the README content. |
| **Streaming output** | As the model generates text, it is streamed to the terminal, reducing latency. |
| **Model selection** | Specify the Ollama model via the `-model` flag. |
| **Custom prompt file** | Override the default `Prompt.md` path with `-promptFile`. |
| **Logging** | Helpful debug output to stderr. |
| **License** | MIT licensed. |

---

## Prerequisites

| Component | Minimum Version | Notes |
|-----------|-----------------|-------|
| **Go** | 1.23.0 | The `go.mod` file pins this version. |
| **Ollama** | 0.1.0+ | Must be running locally on port 11434. |
| **Repository** | Any | The tool will read all files relative to the current working directory. |
| **.gitignore** | Optional | If missing, the tool will process all files. |

---

## Installation

```bash
git clone https://github.com/your-org/ollama-readme-generator.git
cd ollama-readme-generator
go mod download
```

---

## Usage

```bash
# Basic usage
go run main.go -model=gpt-oss > README.md

# Specify a different prompt file
go run main.go -promptFile=MyPrompt.md -model=phi-3 > README.md
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-model` | `gpt-oss` | Name of the Ollama model to query. |
| `-promptFile` | `Prompt.md` | Path to the prompt file that will be prefixed to the file contents. |

The command outputs the generated README to stdout; redirect it to a file as shown above.

---

## Architecture

### Main package

`main.go` orchestrates the workflow:

1. **Flag parsing** – determines model name and prompt file path.  
2. **Prompt reading** – reads the entire content of the specified prompt file.  
3. **File gathering** – calls `lib.ReadData()` to get a string that contains every file’s name and contents.  
4. **Ollama query** – passes the combined string to `lib.AskOllama()`.  
5. **Exit** – exits with status 0 on success, 1 on failure.

### File processing

Implemented in `lib/file-process.go`:

- **`parseGitIgnore`** reads `.gitignore` and returns a slice of rules.  
- **`isIgnored`** applies each rule to a file path.  
- **`ReadData`** walks the current directory, skips ignored files, and concatenates each file’s name and contents into a single string.

The walk is performed with `filepath.WalkDir`, which yields a `DirEntry` for each file or directory. The code respects directory-level ignores by returning `filepath.SkipDir`.

### Ollama client

Implemented in `lib/ollama-sender.go`:

- **`OllamaRequest`** – struct for the JSON request body.  
- **`OllamaResponse`** – struct for individual JSON response chunks.  
- **`AskOllama`** – constructs the HTTP POST, streams the response via a `json.Decoder`, and prints each chunk’s `Response` field. The loop terminates when `Done` is true.

The request uses `Stream: true` so the model streams tokens back incrementally, which reduces perceived latency.

---

## Functions

### `parseGitIgnore`

```go
func parseGitIgnore(path string) ([]string, error)
```

- **Purpose** – Load `.gitignore` rules into memory.  
- **Algorithm** – Opens the file, scans line by line, trims whitespace, discards empty lines and comments, and normalizes each rule by trimming leading/trailing slashes. The `.git` directory is added by default.  
- **Return** – A slice of strings (`rules`) and an error if the file cannot be opened.

### `isIgnored`

```go
func isIgnored(name string, rules []string) bool
```

- **Purpose** – Determine whether a given path should be excluded.  
- **Algorithm** – Iterates over all rules, attempting a `filepath.Match`. If the rule matches the name or the name contains the rule, the function returns `true`.  
- **Return** – `true` if the file is ignored, `false` otherwise.

### `ReadData`

```go
func ReadData() string
```

- **Purpose** – Aggregate all file data into a single string.  
- **Algorithm** –  
  1. Get current working directory.  
  2. Load ignore rules.  
  3. Walk the tree with `filepath.WalkDir`.  
  4. Skip ignored paths (`filepath.SkipDir` for directories, `nil` for files).  
  5. For each non‑ignored file, read its contents and append a header `FileName: <path>` followed by `Data:` and the file content.  
- **Return** – A concatenated string containing the file metadata and contents.

### `AskOllama`

```go
func AskOllama(modelName string, prompt string) error
```

- **Purpose** – Send a prompt to an Ollama model and stream the response.  
- **Algorithm** –  
  1. Marshal `OllamaRequest` into JSON.  
  2. POST to `http://localhost:11434/api/generate`.  
  3. Use `json.NewDecoder` to read the response stream.  
  4. For each decoded `OllamaResponse`, print `part.Response`.  
  5. Break when `part.Done` is true.  
- **Return** – `nil` on success, otherwise the error encountered.

---

## Logging & Debug

- The application uses the standard `log` package.  
- Errors during prompt reading or directory traversal are logged and cause the program to exit.  
- Informational messages (e.g., “Prompt okundu!”) are printed to help the user follow the flow.  
- The file processing routine logs each file as it is read, which can be helpful when diagnosing missing files.

---

## FAQ

**Q: Why does the tool ignore `.git`?**  
A: The default rule added by `parseGitIgnore` guarantees that the `.git` directory itself is never processed, preventing the inclusion of repository metadata.

**Q: How does the tool handle large repositories?**  
A: All file contents are concatenated into a single string before being sent to Ollama. If the repository is very large, consider limiting the depth or using a custom prompt that focuses on the most relevant files.

**Q: Can I use a different Ollama endpoint?**  
A: Currently the URL is hard‑coded to `http://localhost:11434/api/generate`. Modify `AskOllama` if you need a different host or port.

**Q: Does it support other models?**  
A: Yes. Pass the model name via the `-model` flag; the request body will contain that name.

---

## Contributing

Feel free to open issues or submit pull requests. The project follows the MIT license, so any improvements or bug fixes are welcome.

---

## License

MIT License. See the bundled `LICENSE` file for details.

---

## Turkish Documentation

### Genel Bakış

Bu araç, bir kod tabanı için otomatik olarak `README.md` dosyası oluşturan küçük bir Go uygulamasıdır.  
Oluşturulan kod dosyalarını okur, `.gitignore` kurallarını dikkate alır, kullanıcı tarafından sağlanan bir prompt ile birleştirir, bir Ollama modeline gönderir ve modelin yanıtını doğrudan stdout’a akış olarak yazar. Çıktıyı `README.md` olarak yeniden yönlendirebilirsiniz.

---

## İçindekiler

- [Genel Bakış](#genel-bakış)
- [Özellikler](#özellikler)
- [Gereksinimler](#gereksinimler)
- [Kurulum](#kurulum)
- [Kullanım](#kullanım)
- [Mimari](#mimari)
  - [Ana paket](#ana-paket)
  - [Dosya işleme](#dosya-işleme)
  - [Ollama istemcisi](#ollama-istemcisi)
- [Fonksiyonlar](#fonksiyonlar)
  - [parseGitIgnore](#parsegitignore)
  - [isIgnored](#isignored)
  - [ReadData](#readdata)
  - [AskOllama](#askollama)
- [Loglama & Hata Ayıklama](#loglama--hata-ayıklama)
- [SSS](#sss)
- [Katkıda Bulunma](#katkıda-bulunma)
- [Lisans](#lisans)

---

## Genel Bakış

Araç aşağıdaki adımları uygular:

1. **Prompt okuma** – bir dosyadan (varsayılan: `Prompt.md`) içerik alır.  
2. **Klasör tarama** – geçerli çalışma dizini içinde tüm dosyaları tarar, `.gitignore` kurallarına uyanları hariç tutar.  
3. **Birleştirme** – prompt ile toplanan dosya içeriklerini birleştirir.  
4. **Ollama’ya gönderme** – `http://localhost:11434/api/generate` adresine HTTP POST gönderir.  
5. **Akışlı çıktı** – Ollama’dan gelen JSON parçalarını alır ve `response` alanını stdout’a yazar.  
6. **Yönlendirme** – stdout’ı dosyaya yönlendirerek tam bir `README.md` oluşturulur.

Kod, ek bir eş zamanlılık veya karmaşık şablonlamadan, net sorumluluk ayrımıyla yazılmıştır.

---

## Özellikler

| Özellik | Açıklama |
|---------|----------|
| **Otomatik dosya toplama** | `.gitignore` kurallarını dikkate alarak klasör taraması. |
| **Prompt enjeksiyonu** | Kullanıcı tanımlı metin, README’nin şeklini belirler. |
| **Akışlı çıktı** | Model çıktısı anlık olarak terminale gönderilir, gecikme azalır. |
| **Model seçimi** | `-model` bayrağıyla Ollama modelini belirleyin. |
| **Özel prompt dosyası** | `-promptFile` bayrağıyla varsayılan `Prompt.md` yerine başka dosya seçin. |
| **Loglama** | Yardımcı debug çıktısı `stderr`’e yazılır. |
| **Lisans** | MIT lisanslı. |

---

## Gereksinimler

| Bileşen | Minimum Sürüm | Notlar |
|---------|--------------|--------|
| **Go** | 1.23.0 | `go.mod` bu sürümü belirtiyor. |
| **Ollama** | 0.1.0+ | Yerel sunucu olarak 11434 portunda çalışıyor. |
| **Klasör** | Herhangi | Araç, geçerli çalışma dizini ile çalışır. |
| **.gitignore** | Opsiyonel | Yoksa tüm dosyalar işlenir. |

---

## Kurulum

```bash
git clone https://github.com/your-org/ollama-readme-generator.git
cd ollama-readme-generator
go mod download
```

---

## Kullanım

```bash
# Temel kullanım
go run main.go -model=gpt-oss > README.md

# Farklı prompt dosyası belirtme
go run main.go -promptFile=MyPrompt.md -model=phi-3 > README.md
```

### Bayraklar

| Bayrak | Varsayılan | Açıklama |
|--------|------------|----------|
| `-model` | `gpt-oss` | Gönderilecek Ollama modelinin adı. |
| `-promptFile` | `Prompt.md` | Prompt dosyasının yolu. |

Komut, oluşturulan README’yi stdout’a gönderir; örnek gibi bir dosyaya yönlendirme yapabilirsiniz.

---

## Mimari

### Ana paket

`main.go` iş akışını yönetir:

1. **Bayrak ayrıştırma** – model adı ve prompt dosyası belirler.  
2. **Prompt okuma** – belirtilen dosyanın tam içeriğini alır.  
3. **Dosya toplama** – `lib.ReadData()` ile dosya içeriği toplanır.  
4. **Ollama sorgusu** – birleşik metni `lib.AskOllama()` ile gönderir.  
5. **Çıkış** – başarıda 0, hata durumunda 1 ile çıkar.

### Dosya işleme

`lib/file-process.go` içinde:

- **`parseGitIgnore`** – `.gitignore` dosyasını okur ve kuralları döndürür.  
- **`isIgnored`** – verilen yolun kurallara uyup uymadığını kontrol eder.  
- **`ReadData`** – `filepath.WalkDir` ile dizini dolaşır, filtre uygular, dosya içeriğini okur ve birleştirir.

Dizin bazlı hariç tutma için `filepath.SkipDir` kullanılır.

### Ollama istemcisi

`lib/ollama-sender.go` içinde:

- **`OllamaRequest`** – JSON istek gövdesi.  
- **`OllamaResponse`** – API’nin her parçayı döndürdüğü yapı.  
- **`AskOllama`** – POST isteği gönderir, yanıtı `json.Decoder` ile parçalar ve `Response` alanını yazdırır. `Done` true olduğunda döngü sonlanır.

---

## Fonksiyonlar

### `parseGitIgnore`

```go
func parseGitIgnore(path string) ([]string, error)
```

- **Amaç** – `.gitignore` kurallarını hafızaya yükler.  
- **Algoritma** – Dosya satır satır okunur, boşluk ve yorum satırları atılır, `/` karakterleri temizlenir; `.git` dizini varsayılan olarak eklenir.  
- **Döndürme** – kural dizisi ve hata (dosya açılamazsa).

### `isIgnored`

```go
func isIgnored(name string, rules []string) bool
```

- **Amaç** – Belirtilen yolun ignore edilip edilmediğini belirler.  
- **Algoritma** – Kurallar üzerinde döner, `filepath.Match` ile eşleştirir; eşleşme ya da içerme durumunda `true`.  
- **Döndürme** – Ignored ise `true`, aksi halde `false`.

### `ReadData`

```go
func ReadData() string
```

- **Amaç** – Tüm dosyaları tek bir string içinde toplar.  
- **Algoritma** –  
  1. Çalışma dizini alır.  
  2. ignore kuralları yükler.  
  3. `filepath.WalkDir` ile dolaşır.  
  4. Ignored yolları atlar (`filepath.SkipDir` veya `nil`).  
  5. Her dosya için başlık ve içerik ekler.  
- **Döndürme** – dosya başlıkları ve içerikleriyle doldurulmuş bir string.

### `AskOllama`

```go
func AskOllama(modelName string, prompt string) error
```

- **Amaç** – Prompt’u Ollama modeline gönderir ve yanıtı akış olarak alır.  
- **Algoritma** –  
  1. `OllamaRequest` JSON’a dönüştürülür.  
  2. `POST` ile `http://localhost:11434/api/generate` adresine gönderilir.  
  3. `json.Decoder` ile gelen parçalar okunur.  
  4. Her parçanın `Response` alanı yazdırılır.  
  5. `Done` true olduğunda döngü sonlanır.  
- **Döndürme** – Başarılı ise `nil`, aksi halde hata.

---

## Loglama & Hata Ayıklama

- Standart `log` paketi kullanılır.  
- Prompt okuma veya dizin tarama hataları loglanır ve program durur.  
- Bilgilendirme mesajları (örn. “Prompt okundu!”) akışı izlemek için yazdırılır.  
- Dosya okuma sırasında her dosya adı loglanır, eksik dosyaların bulunmasına yardımcı olur.

---

## SSS

**Sorun:** `.git` klasörü neden hariç tutuluyor?  
**Cevap:** `parseGitIgnore` fonksiyonu, `.git` klasörünü varsayılan ignore kuralı olarak ekler; bu, yerel Git metadata’nın README’ye dahil edilmesini önler.

**Sorun:** Çok büyük bir repo için performans nasıl?  
**Cevap:** Tüm dosya içeriği tek bir string içinde birleştirilir. Çok büyük repo varsa derinlik sınırlaması veya kritik dosyalara odaklanma önerilir.

**Sorun:** Farklı bir Ollama endpoint’i kullanabilir miyim?  
**Cevap:** Şu anda URL `http://localhost:11434/api/generate` olarak sabittir. Başka bir ana bilgisayar veya port gerekiyorsa `AskOllama` fonksiyonunda değişiklik yapabilirsiniz.

**Sorun:** Başka modeller kullanılabilir mi?  
**Cevap:** Evet. `-model` bayrağı ile istenen model adı gönderilir; istek gövdesi bu adı içerir.

---

## Katkıda Bulunma

Sorun raporlamak veya pull request göndermek için her zaman memnunumuz. MIT lisansı sayesinde herhangi bir iyileştirme veya hata düzeltmesi serbesttir.

---

## Lisans

MIT Lisansı. Detaylar için ekli `LICENSE` dosyasını inceleyin.

---

Bu dosya AI üzerinden otomatik hazırlanmıştır.

## AI Context & Memory

The repository implements a minimal README generator powered by Ollama. The main program reads a user‑defined prompt (from `Prompt.md` by default), recursively walks the current directory to gather all source files while respecting `.gitignore` rules, and concatenates file names and contents into a single string. It then constructs an `OllamaRequest` with the specified model and a `stream: true` flag, POSTs to the local Ollama API (`http://localhost:11434/api/generate`), and streams the JSON chunks back to stdout. The output can be redirected to produce a fully‑formed `README.md`. Key functions include `parseGitIgnore` (load ignore rules), `isIgnored` (apply rules), `ReadData` (aggregate file data), and `AskOllama` (HTTP client and streaming decoder). Logging is done with Go’s `log` package, and errors cause graceful exit. The project is Go 1.23, uses the standard library, and is MIT licensed.
