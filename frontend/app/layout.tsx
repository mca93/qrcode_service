import "./globals.css";
import { ReactNode } from "react";

export const metadata = {
  title: "QRCode Dashboard",
  description: "Dashboard de gestão de QR Codes",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="pt">
      <body className="bg-gray-50 text-gray-900">
        <nav className="bg-white shadow p-4 flex justify-between items-center">
          <h1 className="text-xl font-bold">QR Dashboard</h1>
        </nav>
        <main className="p-6 max-w-6xl mx-auto">{children}</main>
      </body>
    </html>
  );
}
