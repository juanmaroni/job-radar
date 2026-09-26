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
