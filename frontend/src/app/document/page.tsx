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

export default function DocumentConverter() {
  const [file, setFile] = useState<File | null>(null);
  const [format, setFormat] = useState("pdf");
  const [isLoading, setIsLoading] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      // Validate file type (PDF or DOCX)
      const allowedTypes = [
        "application/pdf",
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      ];
      if (!allowedTypes.includes(selectedFile.type)) {
        toast.error("Please upload a PDF or DOCX file");
        setFile(null);
        return;
      }
      setFile(selectedFile);
    }
  };

  const handleConvert = async () => {
    if (!file) {
      toast.error("Please upload a document");
      return;
    }

    setIsLoading(true);
    const toastId = toast.loading(
      `Converting to ${format.toUpperCase()}... This may take a moment.`,
      {
        description: "File will save to your default Downloads folder.",
      }
    );

    try {
      const API_URL =
        process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const formData = new FormData();
      formData.append("file", file);
      formData.append("format", format);

      const res = await fetch(`${API_URL}/document/convert`, {
        method: "POST",
        headers: {
          Authorization: localStorage.getItem("token") || "",
        },
        body: formData,
      });

      if (!res.ok) {
        const error = await res.text();
        throw new Error(error || "Failed to convert document");
      }

      const blob = await res.blob();
      const contentDisposition = res.headers.get("Content-Disposition");
      const filenameMatch = contentDisposition?.match(/filename="(.+)"/);
      const filename = filenameMatch
        ? filenameMatch[1]
        : `converted_document.${format}`;

      const downloadUrl = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = downloadUrl;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(downloadUrl);

      toast.success(`Downloaded ${filename}`, {
        id: toastId,
        description:
          "Check your Downloads folder or browser settings to locate the file.",
      });
    } catch (error) {
      const err = error as Error;
      toast.error(err.message || "Failed to convert document", {
        id: toastId,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4 text-center">
        Document Converter
      </h1>
      <div className="flex flex-col gap-4 max-w-md mx-auto">
        <Input
          type="file"
          accept=".pdf,.docx"
          onChange={handleFileChange}
          disabled={isLoading}
        />
        <Select value={format} onValueChange={setFormat} disabled={isLoading}>
          <SelectTrigger>
            <SelectValue placeholder="Select format" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="pdf">PDF</SelectItem>
            <SelectItem value="docx">DOCX</SelectItem>
          </SelectContent>
        </Select>
        <Button onClick={handleConvert} disabled={isLoading || !file}>
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Converting...
            </>
          ) : (
            "Convert"
          )}
        </Button>
      </div>
    </div>
  );
}
