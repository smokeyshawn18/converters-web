"use client";

import { useState } from "react";

import Link from "next/link";
import { Menu } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ModeToggle } from "@/components/mode-toggle";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";

export default function Navbar() {
  const [isOpen, setIsOpen] = useState(false);

  const navItems = [
    { href: "/url", label: "URL Shortener" },
    { href: "/image", label: "Image Converter" },
    { href: "/document", label: "Document Converter" },
  ];

  return (
    <header className="border-b p-4 flex items-center justify-between">
      {/* Left: Logo */}
      <h1 className="text-xl font-bold">
        <Link href="/">Utility App</Link>
      </h1>

      {/* Center: Desktop nav links */}
      <nav className="hidden md:flex space-x-8 flex-grow justify-center">
        {navItems.map((item) => (
          <Link key={item.href} href={item.href}>
            <Button variant="ghost">{item.label}</Button>
          </Link>
        ))}
      </nav>

      {/* Right: desktop mode toggle & mobile menu */}
      <div className="flex items-center space-x-2">
        <div className="hidden md:block">
          <ModeToggle />
        </div>

        {/* Mobile menu button */}
        <Sheet open={isOpen} onOpenChange={setIsOpen}>
          <SheetTrigger asChild className="md:hidden">
            <Button variant="ghost" size="icon" aria-label="Menu">
              <Menu className="h-6 w-6" />
            </Button>
          </SheetTrigger>

          <SheetContent side="right" className="w-60">
            <nav className="flex flex-col space-y-4 mt-4">
              {navItems.map((item) => (
                <Link
                  href={item.href}
                  key={item.href}
                  onClick={() => setIsOpen(false)}
                >
                  <Button variant="ghost" className="w-full justify-start">
                    {item.label}
                  </Button>
                </Link>
              ))}
              <div className="pt-4">
                <ModeToggle />
              </div>
            </nav>
          </SheetContent>
        </Sheet>
      </div>
    </header>
  );
}
