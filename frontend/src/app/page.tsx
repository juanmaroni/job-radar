import { Job } from "@/types/job";
import { PageTitle } from "@/components/PageTitle";
import { JobsPanel } from "@/components/JobsPanel";

async function getJobs(): Promise<Job[]> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_JOBS_SERVICE_URL}/jobs`, { cache: "no-store" });
  
  if (!res.ok) {
    throw new Error("Failed to fetch jobs");
  }

  return res.json();
}

export default async function Home() {
  const jobs = await getJobs();

  return (
    <main className="p-8 h-screen">
      <PageTitle title="Jobs" />
      <JobsPanel jobs={jobs} />
    </main>
  );
}
