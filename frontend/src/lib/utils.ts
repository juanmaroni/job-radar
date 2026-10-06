import { theme } from "./theme"

export function scoreColor(score: number | null): string {
  if (score === null) {
    return theme.scoreNone
  } else if (score <= 0) {
    return theme.scoreZero
  } else if (score === 1) {
    return theme.scoreLow
  } else if (score === 2) {
    return theme.scoreMedium
  } else {
    return theme.scoreHigh
  }
}

export function statusColor(status: string): string {
  if (status === "interview") {
    return theme.statusInterview
  } else if (status === "offer") {
    return theme.statusOffer
  } else if (status === "rejected") {
    return theme.statusRejected
  } else {
    return theme.statusApplied
  }
}
