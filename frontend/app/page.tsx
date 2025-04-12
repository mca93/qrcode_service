"use client";
import { useEffect, useState } from "react";
import axios from "axios";

export default function Dashboard() {
  const [qrCount, setQrCount] = useState(0);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await axios.get(`${process.env.NEXT_PUBLIC_API_URL}/qrcodes`, {
          headers: { "X-API-Key": process.env.NEXT_PUBLIC_API_KEY },
        });
        setQrCount(res.data.qr_codes.length);
      } catch (err) {
        console.error("Erro ao carregar QRCodes", err);
      }
    };
    fetchData();
  }, []);

  return (
    <div>
      <h2 className="text-2xl font-semibold mb-4">Dashboard</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white shadow rounded-lg p-4">
          <h3 className="text-sm text-gray-500">Total de QRCodes</h3>
          <p className="text-3xl font-bold">{qrCount}</p>
        </div>
      </div>
    </div>
  );
}
