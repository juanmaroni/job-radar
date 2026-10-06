"use client";

import { Job } from "@/types/job";
import { scoreColor } from "@/lib/utils";

export function JobResumeCard({ job, onSelect }: { job: Job | null; onSelect: (job: Job | null) => void; }) {
  if (job === null) return null;
  
  return (
    <div
      className={`${scoreColor(job.Score)} w-100 h-40 p-5 rounded-md align-middle cursor-pointer`}
      onClick={() => {onSelect(job)}}
    >
      <p className="font-medium">{job.Role}</p>
      <p className="text-sm">{job.Company}</p>
      <p className="text-sm">🎯 {job.Score}</p>
    </div>
  )
}
