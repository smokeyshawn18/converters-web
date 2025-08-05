package handlers

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log/slog"
	"net/http"

	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ImageConvertHandler(pool *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        format := strings.ToLower(r.URL.Query().Get("format"))
        if format == "" {
            format = "jpeg"
        }
        supportedFormats := map[string]bool{"jpeg": true, "png": true, "gif": true}
        if !supportedFormats[format] {
            logger.Warn("Unsupported output format", "format", format)
            http.Error(w, "Unsupported output format", http.StatusBadRequest)
            return
        }

        file, header, err := r.FormFile("file")
        if err != nil {
            logger.Error("Failed to get file", "error", err)
            http.Error(w, "Invalid file upload", http.StatusBadRequest)
            return
        }
        defer file.Close()

        img, imgFormat, err := image.Decode(file)
        if err != nil {
            logger.Error("Failed to decode image", "error", err)
            http.Error(w, "Unsupported image format", http.StatusBadRequest)
            return
        }
        logger.Info("Image decoded", "filename", header.Filename, "format", imgFormat)

        outFilename := fmt.Sprintf("converted_%s.%s", strings.TrimSuffix(header.Filename, "."+imgFormat), format)

        w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", outFilename))
        switch format {
        case "jpeg":
            w.Header().Set("Content-Type", "image/jpeg")
            err = jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
        case "png":
            w.Header().Set("Content-Type", "image/png")
            err = png.Encode(w, img)
        case "gif":
            w.Header().Set("Content-Type", "image/gif")
            err = gif.Encode(w, img, nil)
        default:
            http.Error(w, "Unsupported format", http.StatusBadRequest)
            return
        }

        if err != nil {
            logger.Error("Failed to encode image", "error", err)
            http.Error(w, "Failed to encode image", http.StatusInternalServerError)
            return
        }
     
    }
}
