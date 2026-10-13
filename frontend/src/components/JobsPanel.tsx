"use client";

import { useState } from "react";
import { Job } from "@/types/job";
import { JobResumeCard } from "./JobResumeCard";
import { JobInfoCard } from "./JobInfoCard";
import DOMPurify from "dompurify";

export function JobsPanel({ jobs }: {jobs: Job[]; days: string}) {
  const [selectedJob, setSelectedJob] = useState<Job | null>(null);
  const [searchText, setSearchText] = useState<string>("");

  function handleSearchText(event: React.ChangeEvent<HTMLInputElement>) {
    setSearchText(event.target.value);
  }

  function filterJobs(): Job[] {
    if (jobs === null) return [];

    if (searchText === "") return jobs;

    const searchTextLowercase = searchText.toLocaleLowerCase();
    
    return jobs.filter(
      (job) => job.Role.toLowerCase().includes(searchTextLowercase)
      || job.Company.toLowerCase().includes(searchTextLowercase)
      || DOMPurify.sanitize(job.Description, {
          ALLOWED_TAGS: [],
          ALLOWED_ATTR: []
          }).toLowerCase().includes(searchTextLowercase)
    );
  }

  const filteredJobs = filterJobs();

  return (
    <div>
      <div className="flex gap-4">
        <input
          className="mb-4 p-2 border-2 rounded-md"
          type="text"
          placeholder="Search in text"
          value={searchText}
          onChange={handleSearchText}
        />
      </div>
      <div className="flex gap-2">
        <div className="flex flex-wrap gap-4 justify-flex-start flex-1">
          {filteredJobs.map((job) => (
            <div key={job.Id}>
              <JobResumeCard job={job} onSelect={setSelectedJob} />
          </div>
        ))}
      </div>
        <div className="w-1/3">
          {selectedJob && 
          <JobInfoCard job={selectedJob} />
          }
        </div>
      </div>
    </div>
  )
}
