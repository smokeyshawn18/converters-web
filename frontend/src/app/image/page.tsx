"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";

export default function ImageConverter() {
  const [file, setFile] = useState<File | null>(null);
  const [format, setFormat] = useState("jpeg");
  const [isLoading, setIsLoading] = useState(false);

  const handleConvert = async () => {
    if (!file) {
      toast.error("Please select an image to convert");
      return;
    }

    setIsLoading(true);
    const toastId = toast.loading("Converting image...");

    try {
      const API_URL =
        process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const formData = new FormData();
      formData.append("file", file);

      const res = await fetch(`${API_URL}/image/convert?format=${format}`, {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText || "Conversion failed");
      }

      const blob = await res.blob();
      const outFilename = `converted_image.${format}`;

      const downloadUrl = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = downloadUrl;
      a.download = outFilename;
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(downloadUrl);

      toast.success("Image converted successfully", { id: toastId });
    } catch (err) {
      toast.error((err as Error).message || "Conversion failed", {
        id: toastId,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4 max-w-md">
      <h1 className="text-2xl font-bold mb-4">Image Converter</h1>
      <div className="flex flex-col gap-4">
        <Input
          type="file"
          accept="image/png, image/jpeg, image/gif"
          onChange={(e) => {
            if (e.target.files && e.target.files.length > 0) {
              setFile(e.target.files[0]);
            } else {
              setFile(null);
            }
          }}
          disabled={isLoading}
        />
        <Select value={format} onValueChange={setFormat} disabled={isLoading}>
          <SelectTrigger>
            <SelectValue placeholder="Select format" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="jpeg">JPEG</SelectItem>
            <SelectItem value="png">PNG</SelectItem>
            <SelectItem value="gif">GIF</SelectItem>
          </SelectContent>
        </Select>
        <Button onClick={handleConvert} disabled={isLoading}>
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Converting...
            </>
          ) : (
            "Convert Image"
          )}
        </Button>
      </div>
    </div>
  );
}
