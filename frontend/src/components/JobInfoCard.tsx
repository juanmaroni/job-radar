"use client";

import { Job } from "@/types/job";
import { ApplyButton } from "./ApplyButton";
import { HighlightText } from "./HighlightText";
import { highlightHtml } from "@/lib/highlight";
import { theme } from "@/lib/theme";
import DOMPurify from "dompurify";

export function JobInfoCard({ job, query }: { job: Job | null; query: string }) {
  if (job === null) return null;

  return (
    <div className={`border-2 border-white bg-zinc-800 p-5 rounded-md align-middle fixed h-[80%] overflow-y-auto mr-8`}>
      <p className="text-2xl"><HighlightText text={job.Role} query={query} /></p>
      <p className="text-xl"><HighlightText text={job.Company} query={query} /></p>
      <div className="flex flex-col gap-1 mt-10 text-md h-40">
        <p className="flex">
          <span className={`${theme.spanJobsInfoElems}`}>🎯</span>
          <span>{job.Score}</span>
        </p>
        <p className="flex" hidden={job.Location === ""}>
          <span className={`${theme.spanJobsInfoElems}`}>📍</span>
          <span>{job.Location}</span>
        </p>
        <p className="flex" hidden={job.PublishedAt === ""}>
          <span className={`${theme.spanJobsInfoElems}`}>📅</span>
          <span>{job.PublishedAt.split("T")[0]}</span>
        </p>
        <p className="flex" hidden={job.SalaryMin === null && job.SalaryMax === null}>
          <span className={`${theme.spanJobsInfoElems}`}>💰</span>
          <span>{job.SalaryMin} - {job.SalaryMax}</span>
        </p>
      </div>
      <div
        className={`text-sm **:${theme.text}!`}
        dangerouslySetInnerHTML={{ __html: highlightHtml(DOMPurify.sanitize(job.Description), query) }}
      />
      <div className="flex justify-center gap-2 m-10">
        <button className={`${theme.buttonJobLink} ${theme.buttonJobDescription}`}>
          <a href={job.OfferUrl} target="_blank">Apply</a>
        </button>
        <ApplyButton job={job} />
      </div>
    </div>
  )
}
