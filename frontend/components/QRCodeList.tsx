"use client";
import { useEffect, useState } from "react";
import axios from "axios";

export default function QRCodeList() {
  const [codes, setCodes] = useState([]);

  useEffect(() => {
    axios.get(`${process.env.NEXT_PUBLIC_API_URL}/qrcodes`, {
      headers: { "X-API-Key": process.env.NEXT_PUBLIC_API_KEY },
    })
    .then(res => setCodes(res.data.qr_codes))
    .catch(err => console.error(err));
  }, []);

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">QRCodes</h2>
      <table className="min-w-full bg-white rounded shadow">
        <thead>
          <tr>
            <th className="p-2 text-left">ID</th>
            <th className="p-2 text-left">Status</th>
          </tr>
        </thead>
        <tbody>
          {codes.map((code: any) => (
            <tr key={code.id}>
              <td className="p-2">{code.id}</td>
              <td className="p-2">{code.status}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
