import { Application } from "@/types/application";
import { StatusButtons } from "@/components/StatusButton";

async function getApplications(): Promise<Application[]> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_TRACKER_SERVICE_URL}/applications`, { cache: "no-store" });
  
  if (!res.ok) {
    throw new Error("Failed to fetch applications");
  }

  return res.json();
}

function groupByStatus(applications: Application[]): Record<string, Application[]> {
  return applications.reduce((groups, app) => {
    const status = app.Status;
    
    if (!groups[status]) {
      groups[status] = [];
    }
    
    groups[status].push(app);

    return groups;
  }, {} as Record<string, Application[]>);
}

export default async function Home() {
  const applications = await getApplications();
  
  const STATUSES = ["applied", "interview", "offer", "rejected"];
  const grouped = groupByStatus(applications);

  return (
    <main>
      <h1>Applications</h1>
      <div style={{ display: "flex", gap: "1rem" }}>
      {STATUSES.map((status) => (
        <div key={status}>
          <h2>{status}</h2>
          {(grouped[status] ?? []).map((app) => (
            <div key={app.Id}>
              <p>{app.Company} — {app.Role}</p>
              <StatusButtons application={app} />
            </div>
          ))}
        </div>
      ))}
      </div>
    </main>
  );
}
