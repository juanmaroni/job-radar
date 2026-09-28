"use client";

import { useDraggable } from "@dnd-kit/core";
import { Application } from "@/types/application";
import { StatusButtons } from "./StatusButton";

export function KanbanCard({
  application,
  onStatusChange,
}: {
  application: Application;
  onStatusChange: (id: number, status: string) => void;
}) {
  const { attributes, listeners, setNodeRef, transform } = useDraggable({
    id: application.Id.toString(),
  });

  const style = transform
    ? { transform: `translate(${transform.x}px, ${transform.y}px)` }
    : undefined;

  return (
    <div ref={setNodeRef} style={style} className="bg-white rounded p-3 shadow">
      <div {...listeners} {...attributes} className="cursor-grab">
        <p className="font-medium text-black">{application.Role}</p>
        <p className="text-sm text-gray-500">{application.Company}</p>
      </div>
      <StatusButtons application={application} onStatusChange={onStatusChange} />
    </div>
  );
}
