package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

  func DocConvertHandler(pool *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
  	return func(w http.ResponseWriter, r *http.Request) {
  		if r.Method != http.MethodPost {
  			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  			return
  		}

  		file, _, err := r.FormFile("file")
  		if err != nil {
  			logger.Error("Failed to get file", "error", err)
  			http.Error(w, "Invalid file", http.StatusBadRequest)
  			return
  		}
  		defer file.Close()

  		// Save uploaded file
  		filename := uuid.New().String() + ".pdf"
  		outFile, err := os.Create(filename)
  		if err != nil {
  			logger.Error("Failed to create file", "error", err)
  			http.Error(w, "Failed to save file", http.StatusInternalServerError)
  			return
  		}
  		if _, err := io.Copy(outFile, file); err != nil {
  			logger.Error("Failed to save file", "error", err)
  			http.Error(w, "Failed to save file", http.StatusInternalServerError)
  			return
  		}
  		outFile.Close()
  		defer os.Remove(filename)

  		// Convert to DOCX using LibreOffice
  		outputDir := "converted"
  		if err := os.MkdirAll(outputDir, 0755); err != nil {
  			logger.Error("Failed to create output directory", "error", err)
  			http.Error(w, "Conversion failed", http.StatusInternalServerError)
  			return
  		}
  		cmd := exec.Command("libreoffice", "--headless", "--convert-to", "docx", filename, "--outdir", outputDir)
  		if err := cmd.Run(); err != nil {
  			logger.Error("Document conversion failed", "error", err)
  			http.Error(w, "Conversion failed", http.StatusInternalServerError)
  			return
  		}

  		// Serve converted file
  		convertedFile := filepath.Join(outputDir, filename[:len(filename)-4]+".docx")
  		defer os.Remove(convertedFile)
  		w.Header().Set("Content-Disposition", "attachment; filename=converted_document.docx")
  		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
  		http.ServeFile(w, r, convertedFile)
  	}
  }