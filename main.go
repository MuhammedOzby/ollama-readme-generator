package main

import (
	"flag"
	"fmt"
	"log"
	"ollama-readme-generator/lib"
	"os"
	"path/filepath"
)

func readPromptFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("prompt dosyası okunamadı (%s): %v", filename, err)
	}
	return string(content), nil
}

func main() {
	appDir, err := os.Executable()
	if err != nil {
		fmt.Println("Executable yolu alınamadı:", err)
		return
	}
	modelPtr := flag.String("model", "gpt-oss", "Kullanılacak Ollama modeli")
	promptFilePtr := flag.String("promptfile", filepath.Join(filepath.Dir(appDir), "Prompt.md"), "Code öncesi girilecek olan propmt metinin dosya konumu.")
	flag.Parse()
	prompt, err := readPromptFile(*promptFilePtr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Model: " + *modelPtr)
	log.Println("Promp file location:" + *promptFilePtr)
	log.Println("Prompt okundu! Dosya içeriği alınıyor.")
	filesData := lib.ReadData()
	log.Println("Dosyalar okundu! ollama sorgusu alınıyor.")
	log.Println("Örnek kullanım: go run main.go -model=gpt-oss > README.md")
	log.Println("-----------------------------------------")
	lib.AskOllama(*modelPtr, fmt.Sprintf("%s\n%s", prompt, filesData))
	os.Exit(0)
}
