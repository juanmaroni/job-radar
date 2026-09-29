export function scoreColor(score: number | null): string {
  if (score === null) {
    return "bg-gray-400"
  } else if (score === 0) {
    return "bg-red-400"
  } else if (score <= 4) {
    return "bg-orange-300"
  } else if (score <= 8) {
    return "bg-amber-400"
  } else {
    return "bg-green-400"
  }
}

export function statusColor(status: string): string {
  if (status === "interview") {
    return "bg-amber-400"
  } else if (status === "offer") {
    return "bg-green-400"
  } else if (status === "rejected") {
    return "bg-red-400"
  } else {
    return "bg-gray-400"
  }
}
