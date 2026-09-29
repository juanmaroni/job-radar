"use client";

import { useDraggable } from "@dnd-kit/core";
import { Application } from "@/types/application";
import { StatusButtons } from "./StatusButton";
import { scoreColor } from "@/lib/utils";

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
    <div ref={setNodeRef} style={style} className={`${scoreColor(application.MatchScore)} rounded p-3`}>
      <div {...listeners} {...attributes} className="cursor-grab">
        <p className="font-medium text-black">{application.Role}</p>
        <p className="text-sm text-black">{application.Company}</p>
        <p className="text-sm text-black">🎯 {application.MatchScore}</p>
      </div>
      <StatusButtons application={application} onStatusChange={onStatusChange} />
    </div>
  );
}
