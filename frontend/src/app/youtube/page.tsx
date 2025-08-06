"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";

export default function YouTubeToMP3() {
  const [url, setUrl] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleDownload = async () => {
    if (!url.trim()) {
      toast.error("Please enter a YouTube URL");
      return;
    }

    setIsLoading(true);
    const toastId = toast.loading("Processing your download...");

    try {
      const API_URL = process.env.BACKEND_URL || "http://localhost:8080";

      const response = await fetch(
        `${API_URL}/youtube/mp3?url=${encodeURIComponent(url)}`,
        {
          method: "POST",
        }
      );

      if (!response.ok) {
        const errText = await response.text();
        throw new Error(errText || "Failed to download");
      }

      const blob = await response.blob();
      const downloadUrl = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = downloadUrl;
      a.download = "audio.mp3";
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(downloadUrl);

      toast.success("Download started", { id: toastId });
    } catch (e) {
      toast.error((e as Error).message || "Download failed", { id: toastId });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4 max-w-md">
      <h1 className="text-2xl font-bold mb-4">YouTube to MP3 Downloader</h1>
      <Input
        type="url"
        placeholder="Enter YouTube URL"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        disabled={isLoading}
      />
      <Button
        onClick={handleDownload}
        disabled={isLoading}
        className="mt-4 w-full"
      >
        {isLoading ? (
          <>
            <Loader2 className="mr-2 w-5 h-5 animate-spin" /> Processing...
          </>
        ) : (
          "Download MP3"
        )}
      </Button>
    </div>
  );
}
