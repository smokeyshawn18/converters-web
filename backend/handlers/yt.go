package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"
)

func YoutubeToMP3Handler(logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        url := r.URL.Query().Get("url")
        if url == "" {
            http.Error(w, "Missing url parameter", http.StatusBadRequest)
            return
        }

        // Sanitize/validate URL if needed
        if !strings.Contains(url, "youtube.com") && !strings.Contains(url, "youtu.be") {
            http.Error(w, "Invalid YouTube URL", http.StatusBadRequest)
            return
        }

        // Set response headers for file download
        w.Header().Set("Content-Disposition", `attachment; filename="audio.mp3"`)
        w.Header().Set("Content-Type", "audio/mpeg")

        // Run yt-dlp to extract audio only as mp3, stream output to stdout
        cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "-o", "-", url)

        // Pipe stdout to the response writer
        stdout, err := cmd.StdoutPipe()
        if err != nil {
            logger.Error("Failed to get stdout pipe", "error", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        defer stdout.Close()

        // Start the command
        if err := cmd.Start(); err != nil {
            logger.Error("Failed to start yt-dlp", "error", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        // Stream the output to client
        _, err = fmt.Fprint(w, "")
        if err != nil {
            return
        }
        _, err = io.Copy(w, stdout)
        if err != nil {
            logger.Error("Failed to stream audio", "error", err)
        }

        // Wait for yt-dlp to finish
        if err := cmd.Wait(); err != nil {
            logger.Error("yt-dlp command failed", "error", err)
        }
    }
}
