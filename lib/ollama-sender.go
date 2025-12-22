package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OllamaRequest API'ye gönderilecek veri yapısı
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse API'den gelecek her bir parça için yapı
type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func AskOllama(modelName string, prompt string) error {
	url := "http://localhost:11434/api/generate"

	// İsteği oluştur
	reqBody := OllamaRequest{
		Model:  modelName,
		Prompt: prompt,
		Stream: true, // Yanıtı parça parça almak için
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// HTTP POST isteği gönder
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ollama'ya ulaşılamadı: %v", err)
	}
	defer resp.Body.Close()

	// Yanıtı stream olarak oku
	decoder := json.NewDecoder(resp.Body)
	for {
		var part OllamaResponse
		if err := decoder.Decode(&part); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// Gelen metin parçasını ekrana yazdır
		fmt.Print(part.Response)

		if part.Done {
			break
		}
	}

	fmt.Println() // Sonunda bir alt satıra geç
	return nil
}
