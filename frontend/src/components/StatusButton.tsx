"use client";

import { useState } from "react";
import { Application } from "@/types/application";

const STATUSES = ["applied", "interview", "offer", "rejected"];

export function StatusButtons({ application }: { application: Application }) {
  const [currentStatus, setCurrentStatus] = useState(application.Status);
  const [loading, setLoading] = useState(false);

  async function handleClick(newStatus: string) {
    setLoading(true);

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_TRACKER_SERVICE_URL}/applications/${application.Id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: newStatus }),
      });

      if (!res.ok) {
        throw new Error("Failed to update status");
      }

      setCurrentStatus(newStatus);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      {STATUSES.filter((s) => s !== currentStatus).map((s) => (
        <button key={s} onClick={() => handleClick(s)} disabled={loading} className="text-xs px-2 py-1 bg-blue-500 text-white rounded m-2 cursor-pointer">
          Move to {s}
        </button>
      ))}
    </div>
  );
}
