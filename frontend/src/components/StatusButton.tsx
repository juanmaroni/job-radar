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
        <button key={s} onClick={() => handleClick(s)} disabled={loading}>
          Mover a {s}
        </button>
      ))}
    </div>
  );
}
