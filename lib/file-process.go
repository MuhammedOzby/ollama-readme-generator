package lib

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// .gitignore kurallarını temizleyen ve listeye alan fonksiyon
func parseGitIgnore(path string) ([]string, error) {
	var rules []string = []string{".git"}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Boş satırları ve yorumları atla
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Go'nun filepath.Match yapısına uygun hale getir (basit düzeyde)
		line = strings.Trim(line, "/")
		rules = append(rules, line)
	}
	return rules, scanner.Err()
}

func isIgnored(name string, rules []string) bool {
	for _, rule := range rules {
		// Basit eşleşme kontrolü
		match, _ := filepath.Match(rule, name)
		if match || strings.Contains(name, rule) {
			return true
		}
	}
	return false
}

func ReadData() string {
	root, _ := os.Getwd()
	var filesDatas string

	// 1. Gitignore kurallarını yükle
	ignoreRules, err := parseGitIgnore(root + "/.gitignore")
	if err != nil {
		log.Println(".gitignore bulunamadı veya okunamadı, filtreleme yapılmayacak.")
	}

	// 2. Tarama
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Bağıl yolu (relative path) hesapla ki kural kontrolü kolay olsun
		relPath, _ := filepath.Rel(root, path)

		// .gitignore kontrolü
		if relPath != "." && isIgnored(relPath, ignoreRules) {
			if d.IsDir() {
				return filepath.SkipDir // Dizinse içine hiç girme
			}
			return nil // Dosyaysa işlemeden geç
		}

		// Kod dosyası okuma mantığı
		if !d.IsDir() {
			log.Printf("Okunuyor: %s\n", relPath)
			// content, _ := os.ReadFile(path) ...
			content, _ := os.ReadFile(path)
			filesDatas += fmt.Sprintf("FileName: %s\n", path)
			filesDatas += fmt.Sprintf("Data:\n%s\n\n", string(content))
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Hata: %v\n", err)
	}

	return filesDatas
}
