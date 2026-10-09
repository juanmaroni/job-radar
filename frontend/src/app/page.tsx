import { Job } from "@/types/job";
import { PageTitle } from "@/components/PageTitle";
import { DaysSelect } from "@/components/DaysSelect";
import { JobsPanel } from "@/components/JobsPanel";
import { theme } from "@/lib/theme";

async function getJobs(days: string): Promise<Job[]> {
  let queryDays = "";

  if (days !== "all") {
    queryDays = `?days=${days}`;
  }

  const res = await fetch(`${process.env.NEXT_PUBLIC_JOBS_SERVICE_URL}/jobs${queryDays}`, { cache: "no-store" });
  
  if (!res.ok) {
    throw new Error("Failed to fetch jobs");
  }

  return res.json();
}

export default async function Home({ searchParams }: {searchParams: Promise<{ days?: string }>}) {
  let days = (await searchParams).days;

  if (days === undefined || (days !== "all" && +days <= 0)) {
    days = "30";
  }

  const jobs = await getJobs(days);

  return (
    <main className={`p-8 ${theme.background} ${theme.text}`}>
      <PageTitle title="Jobs" />
      <DaysSelect current={days} />
      {jobs !== null && <JobsPanel jobs={jobs} key={days} />}
    </main>
  );
}
