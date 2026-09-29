"use client";

import { useState } from "react";
import { Application } from "@/types/application";
import { updateApplicationStatus } from "@/lib/applications";

const STATUSES = ["applied", "interview", "offer", "rejected"];

export function StatusButtons({ application, onStatusChange }: { application: Application; onStatusChange: (id: number, status: string) => void }) {
  const [loading, setLoading] = useState(false);

  async function handleClick(newStatus: string) {
    setLoading(true);

    try {
      await updateApplicationStatus(application.Id, newStatus)
      onStatusChange(application.Id, newStatus);
    } catch(e) {
        console.error(e)
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      {STATUSES.filter((s) => s !== application.Status).map((s) => (
        <button key={s} onClick={() => handleClick(s)} disabled={loading} className="text-xs px-2 py-1 bg-blue-500 text-white rounded m-2 cursor-pointer">
          Move to {s}
        </button>
      ))}
    </div>
  );
}
