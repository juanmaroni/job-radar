"use client";

import { Job } from "@/types/job";
import { ApplyButton } from "@/components/ApplyButton";
import { scoreColor } from "@/lib/utils";

export function JobCard({ job }: { job: Job }) {
  return (
    <div className={`${scoreColor(job.Score)} w-100 h-40 p-5 rounded-md align-middle`}>
      <p className="font-medium text-black">{job.Role}</p>
      <p className="text-sm text-black">{job.Company}</p>
      <p className="text-sm text-black">🎯 {job.Score}</p>
      <ApplyButton job={job} />
    </div>
  )
}
