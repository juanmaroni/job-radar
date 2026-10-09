"use client";

import { useRouter } from "next/navigation";

export function DaysSelect({ current }: {current: string}) {
  const router = useRouter();

  function handleChange(event: React.ChangeEvent<HTMLSelectElement>) {
    const value = event.target.value;
    router.push(`/?days=${value}`);
  }

  return(
    <select className="mb-4 p-2 border-2 rounded-md" value={current} onChange={handleChange}>
      <option value={"1"}>Last day</option>
      <option value={"7"}>Last 7 days</option>
      <option value={"14"}>Last 14 days</option>
      <option value={"30"}>Last 30 days</option>
      <option value={"90"}>Last 90 days</option>
      <option value={"all"}>Anytime</option>
    </select>
  )
}
