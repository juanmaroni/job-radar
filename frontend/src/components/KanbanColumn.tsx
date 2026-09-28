"use client";

import { useDroppable } from "@dnd-kit/core";
import React from "react";

export function KanbanColumn({ status, children}: { status: string, children: React.ReactNode}) {
  const { setNodeRef, isOver } = useDroppable({ id: status });

  return(
    <div ref={setNodeRef} className={`flex-1 bg-gray-100 rounded-lg p-4 ${isOver ? "bg-gray-200" : ""}`}>
    <h2 className="font-semibold mb-3 capitalize text-black">{status}</h2>
    <div className="flex flex-col gap-2">{children}</div>
    </div>
  );
}
