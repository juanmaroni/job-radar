"use client";

import { useState } from "react";
import { Application } from "@/types/application";
import { updateApplicationStatus } from "@/lib/applications";
import { theme } from "@/lib/theme";

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
        <button
          key={s}
          onClick={() => handleClick(s)}
          disabled={loading}
          className={`text-xs px-2 py-1 ${theme.buttonStatusChange} rounded m-2 cursor-pointer`}
        >
          Move to {s}
        </button>
      ))}
    </div>
  );
}
