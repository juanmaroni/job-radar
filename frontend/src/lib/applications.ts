import { Application } from "@/types/application";

export function groupByStatus(applications: Application[]): Record<string, Application[]> {
  return applications.reduce((groups, app) => {
    const status = app.Status;
    
    if (!groups[status]) {
      groups[status] = [];
    }
    
    groups[status].push(app);

    return groups;
  }, {} as Record<string, Application[]>);
}

export async function updateApplicationStatus(id: number, newStatus: string): Promise<void> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_TRACKER_SERVICE_URL}/applications/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ status: newStatus }),
  });

  if (!res.ok) {
    throw new Error("Failed to update status");
  }
}
