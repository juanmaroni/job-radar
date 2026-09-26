import { Application } from "@/types/application";
import { StatusButtons } from "@/components/StatusButton";
import { PageTitle } from "@/components/PageTitle";
import { groupByStatus } from "@/lib/applications";

async function getApplications(): Promise<Application[]> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_TRACKER_SERVICE_URL}/applications`, { cache: "no-store" });
  
  if (!res.ok) {
    throw new Error("Failed to fetch applications");
  }

  return res.json();
}

export default async function Home() {
  const applications = await getApplications();
  
  const STATUSES = ["applied", "interview", "offer", "rejected"];
  const grouped = groupByStatus(applications);

  return (
    <main className="p-8">
      <PageTitle title="Applications" />
      <div className="flex gap-4">
        {STATUSES.map((status) => (
          <div key={status} className="flex-1 bg-gray-100 rounded-lg p-4">
            <h2 className="font-semibold mb-3 capitalize text-black">{status}</h2>
            <div className="flex flex-col gap-2">
              {(grouped[status] ?? []).map((app) => (
                <div key={app.Id} className="bg-white rounded p-3 shadow">
                  <p className="font-medium text-black">{app.Role}</p>
                  <p className="text-sm text-gray-500">{app.Company}</p>
                  <StatusButtons application={app} />
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </main>
  );
}
