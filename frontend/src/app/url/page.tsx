"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { z } from "zod";

const schema = z.object({
  url: z.string().url("Invalid URL"),
});

export default function URLShortener() {
  const [url, setUrl] = useState("");
  const [shortUrl, setShortUrl] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const result = schema.safeParse({ url });
    if (!result.success) {
      toast.error(result.error.issues[0].message); // ✅ CORRECT
      return;
    }

    setIsLoading(true);

    try {
      const API_URL = process.env.BACKEND_URL || "http://localhost:8080";
      const shortened = await toast.promise(
        fetch(`${API_URL}/url/shorten`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: localStorage.getItem("token") || "",
          },
          body: JSON.stringify({ url }),
        }).then(async (res) => {
          if (!res.ok) throw new Error("Failed to shorten URL");
          const data = await res.json();
          setShortUrl(data.shortUrl);
          return data.shortUrl;
        }),
        {
          loading: "Shortening URL...",
          success: (short) => `Short URL created: ${short}`,
          error: (err) => err.message,
        }
      );
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-md mx-auto">
      <h2 className="text-2xl font-bold mb-4">URL Shortener</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          type="text"
          placeholder="Enter URL to shorten"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
        <Button type="submit" disabled={isLoading}>
          {isLoading ? "Shortening..." : "Shorten"}
        </Button>
      </form>

      {shortUrl && (
        <p className="mt-4">
          Short URL:{" "}
          <a
            href={shortUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-500"
          >
            {shortUrl}
          </a>
        </p>
      )}
    </div>
  );
}
