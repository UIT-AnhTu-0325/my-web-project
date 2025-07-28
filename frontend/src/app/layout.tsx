import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Navbar from "@/components/Navbar";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Hotel Commerce - Book Rooms & Shop Products",
  description: "A modern hotel booking and e-commerce platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <Navbar />
        <main className="min-h-screen bg-gray-50">
          {children}
        </main>
        <footer className="bg-gray-800 text-white py-8">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center">
              <h3 className="text-lg font-semibold mb-2">🏨 HotelCommerce</h3>
              <p className="text-gray-300">Your comfort is our priority</p>
              <div className="mt-4 text-sm text-gray-400">
                © 2025 HotelCommerce. All rights reserved.
              </div>
            </div>
          </div>
        </footer>
      </body>
    </html>
  );
}
