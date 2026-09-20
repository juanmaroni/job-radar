import { Job } from "@/types/job";
import { ApplyButton } from "@/components/ApplyButton";

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
    <main>
      <h1>Jobs</h1>
      <ul>
        {jobs.map((job) => (
          <li key={job.Id}>
            {job.Company} — {job.Role} (score: {job.Score})
            <ApplyButton job={job} />
          </li>
        ))}
      </ul>
    </main>
  );
}
