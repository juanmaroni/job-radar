import { Application } from "@/types/application";
import { KanbanBoard } from "@/components/KanbanBoard";
import { PageTitle } from "@/components/PageTitle";

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
    <main className="p-8">
      <PageTitle title="Applications" />
      <KanbanBoard initialApplications={applications} />
    </main>
  );
}
