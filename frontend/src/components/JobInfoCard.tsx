"use client";

import { Job } from "@/types/job";
import { scoreColor } from "@/lib/utils";
import { ApplyButton } from "./ApplyButton";

export function JobInfoCard({ job }: { job: Job | null }) {
  if (job === null) return null;

  return (
    <div className={`${scoreColor(job.Score)} h-full p-5 rounded-md align-middle`}>
      <h2 className="font-medium text-black">{job.Role}</h2>
      <h3 className="text-sm text-black">{job.Company}</h3>
      <p className="text-sm text-black">🎯 {job.Score}</p>
      <p className="text-sm text-black">📍 {job.Location}</p>
      <p className="text-sm text-black">📅 {job.PublishedAt}</p>
      <p className="text-sm text-black">💰 {job.SalaryMin} - {job.SalaryMax}</p>
      <p className="text-sm text-black">{job.Description}</p>
      <div className="flex">
        <button className="text-md text-center px-2 py-1 bg-blue-500  rounded w-50 h-10 m-2 cursor-pointer">
            <a href={job.OfferUrl} target="_blank">Apply</a>
        </button>
        <ApplyButton job={job} />
      </div>
    </div>
  )
}
