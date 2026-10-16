"use client";

import { Job } from "@/types/job";
import { scoreColor } from "@/lib/utils";
import { HighlightText } from "./HighlightText";

export function JobResumeCard({ job, query, onSelect }: { job: Job | null; query: string; onSelect: (job: Job | null) => void; }) {
  if (job === null) return null;
  
  return (
    <div
      className={`${scoreColor(job.Score)} w-100 h-40 p-5 rounded-md align-middle cursor-pointer`}
      onClick={() => {onSelect(job)}}
    >
      <p className="font-medium"><HighlightText text={job.Role} query={query} /></p>
      <p className="text-sm"><HighlightText text={job.Company} query={query} /></p>
      <p className="text-sm">🎯 {job.Score}</p>
    </div>
  )
}
