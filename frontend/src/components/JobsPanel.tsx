"use client";

import { useState } from "react";
import { Job } from "@/types/job";
import { JobResumeCard } from "./JobResumeCard";
import { JobInfoCard } from "./JobInfoCard";

export function JobsPanel({ jobs }: {jobs: Job[]}) {
  const [selectedJob, setSelectedJob] = useState<Job | null>(null)

  return (
    <div className="flex gap-2 screen">
      <div className="flex flex-wrap gap-4 justify-flex-start flex-1 w-1/3 h-full">
        {jobs.map((job) => (
          <div key={job.Id}>
            <JobResumeCard job={job} onSelect={setSelectedJob} />
          </div>
        ))}
      </div>
      <div className="h-full w-1/3 fixed right-0">
        {selectedJob && 
        <JobInfoCard job={selectedJob} />
        }
      </div>
    </div>
  )
}
