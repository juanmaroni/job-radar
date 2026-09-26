"use client";

import { useState } from "react";
import { Application } from "@/types/application";
import { groupByStatus } from "@/lib/applications";
import { StatusButtons } from "./StatusButton";

const STATUSES = ["applied", "interview", "offer", "rejected"];

export function KanbanBoard({ initialApplications }: { initialApplications: Application[] }) {
  const [applications, setApplications] = useState(initialApplications);

  function handleStatusChange(id: number, newStatus: string) {
    setApplications((prev) =>
      prev.map((app) => (app.Id === id ? { ...app, Status: newStatus } : app))
    );
  }

  const grouped = groupByStatus(applications);

  return (
    <div className="flex gap-4">
      {STATUSES.map((status) => (
        <div key={status}>
          <h2>{status}</h2>
          {(grouped[status] ?? []).map((app) => (
            <div key={app.Id} className="bg-white rounded p-3 m-3 shadow">
              <p className="font-medium text-black">{app.Role}</p>
              <p className="text-sm text-gray-500">{app.Company}</p>
              <StatusButtons application={app} onStatusChange={handleStatusChange}  />
            </div>
          ))}
        </div>
      ))}
    </div>
  );
}
