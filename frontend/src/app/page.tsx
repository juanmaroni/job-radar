import { Job } from "@/types/job";
import { ApplyButton } from "@/components/ApplyButton";
import { PageTitle } from "@/components/PageTitle";

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
    <main className="p-8">
      <PageTitle title="Jobs" />
      <div className="flex flex-wrap gap-4 justify-between">
        {jobs.map((job) => (
          <div key={job.Id} className="bg-white rounded p-3 shadow w-100 h-40">
            <p className="font-medium text-black">{job.Role}</p>
            <p className="text-sm text-gray-500">{job.Company}</p>
            <p className="text-sm text-gray-500">Match Score: {job.Score}</p>
            <ApplyButton job={job} />
          </div>
        ))}
      </div>
    </main>
  );
}
