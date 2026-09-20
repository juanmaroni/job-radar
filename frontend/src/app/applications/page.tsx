import { Application } from "@/types/application";

async function getApplications(): Promise<Application[]> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_TRACKER_SERVICE_URL}/applications`, { cache: "no-store" });
  
  if (!res.ok) {
    throw new Error("Failed to fetch applications");
  }

  return res.json();
}

export default async function Home() {
  const applications = await getApplications();

  return (
    <main>
      <h1>Applications</h1>
      <ul>
        {(applications ?? []).map((application) => (
          <li key={application.Id}>
            {application.Company} — {application.Role} (score: {application.MatchScore})
          </li>
        ))}
      </ul>
    </main>
  );
}
