import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Link from "next/link";

export default function Home() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {/* <Card>
        <CardHeader>
          <CardTitle>YouTube Downloader</CardTitle>
          <CardDescription>
            Download YouTube videos in MP3 or MP4 format
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Link href="/youtube">
            <Button>Try Now</Button>
          </Link>
        </CardContent>
      </Card> */}
      <Card>
        <CardHeader>
          <CardTitle>URL Shortener</CardTitle>
          <CardDescription>Shorten long URLs for easy sharing</CardDescription>
        </CardHeader>
        <CardContent>
          <Link href="/url">
            <Button>Try Now</Button>
          </Link>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Image Converter</CardTitle>
          <CardDescription>Convert images to JPEG format</CardDescription>
        </CardHeader>
        <CardContent>
          <Link href="/image">
            <Button>Try Now</Button>
          </Link>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Document Converter</CardTitle>
          <CardDescription>Convert PDFs to DOCX format</CardDescription>
        </CardHeader>
        <CardContent>
          <Link href="/document">
            <Button>Try Now</Button>
          </Link>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Download History</CardTitle>
          <CardDescription>View your YouTube download history</CardDescription>
        </CardHeader>
        <CardContent>
          <Link href="/history">
            <Button>Try Now</Button>
          </Link>
        </CardContent>
      </Card>
    </div>
  );
}
