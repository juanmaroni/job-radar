"use client";

import { Job } from "@/types/job";
import { scoreColor } from "@/lib/utils";
import { ApplyButton } from "./ApplyButton";
import { theme } from "@/lib/theme";
import DOMPurify from "dompurify";

export function JobInfoCard({ job }: { job: Job | null }) {
  if (job === null) return null;

  return (
    <div className={`${scoreColor(job.Score)} p-5 rounded-md align-middle fixed h-[90%] overflow-y-auto mr-8`}>
      <p className="text-2xl">{job.Role}</p>
      <p className="text-xl">{job.Company}</p>
      <div className="p-2 text-md">
        <p>🎯 {job.Score}</p>
        <p>📍 {job.Location}</p>
        <p>📅 {job.PublishedAt}</p>
        <p>💰 {job.SalaryMin} - {job.SalaryMax}</p>
      </div>
      <div className="text-sm" dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(job.Description) }} />
      <div className="flex justify-center gap-2 m-10">
        <button className={`${theme.buttonJobLink} ${theme.buttonJobDescription}`}>
          <a href={job.OfferUrl} target="_blank">Apply</a>
        </button>
        <ApplyButton job={job} />
      </div>
    </div>
  )
}
