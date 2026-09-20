"use client";

import { useState } from "react";
import { Job } from "@/types/job";

export function ApplyButton({ job }: { job: Job }) {
  const [status, setStatus] = useState<"idle" | "loading" | "done" | "error">("idle");

  async function handleClick() {
    setStatus("loading");

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_TRACKER_SERVICE_URL}/applications`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          company: job.Company,
          role: job.Role,
          url: job.OfferUrl,
        }),
      });

      if (!res.ok) {
        throw new Error("Failed to create application");
      }

      setStatus("done");
    } catch {
      setStatus("error");
    }
  }

  return (
    <button onClick={handleClick} disabled={status === "loading" || status === "done"}>
      {status === "idle" && "Aplicar"}
      {status === "loading" && "Aplicando..."}
      {status === "done" && "Aplicado ✓"}
      {status === "error" && "Error, reintentar"}
    </button>
  );
}
