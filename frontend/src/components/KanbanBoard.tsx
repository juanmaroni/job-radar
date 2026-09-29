"use client";

import { useState } from "react";
import { DndContext, DragEndEvent } from "@dnd-kit/core";
import { Application } from "@/types/application";
import { groupByStatus, updateApplicationStatus } from "@/lib/applications";
import { KanbanColumn } from "./KanbanColumn";
import { KanbanCard } from "./KanbanCard";

const STATUSES = ["applied", "interview", "offer", "rejected"];

export function KanbanBoard({ initialApplications }: { initialApplications: Application[] }) {
  const [applications, setApplications] = useState(initialApplications);

  function handleStatusChange(id: number, newStatus: string) {
    setApplications((prev) =>
      prev.map((app) => (app.Id === id ? { ...app, Status: newStatus } : app))
    );
  }

  async function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;

    if (!over) return;

    const id = Number(active.id);
    const newStatus = over.id.toString();

    if (applications.find(app => app.Id === id && app.Status === newStatus)) return;

    try {
      await updateApplicationStatus(id, newStatus);
    } catch(e) {
      console.error(e)
      return
    }

    handleStatusChange(id, newStatus);
  }

  const grouped = groupByStatus(applications);

  return (
    <DndContext onDragEnd={handleDragEnd}>
      <div className="flex gap-4">
        {STATUSES.map((status) => (
          <KanbanColumn key={status} status={status}>
            {(grouped[status] ?? []).map((app) => (
              <KanbanCard key={app.Id} application={app} onStatusChange={handleStatusChange} />
            ))}
          </KanbanColumn>
        ))}
      </div>
    </DndContext>
  );
}
