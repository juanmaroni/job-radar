export interface Job {
  Id: number;
  Source: string;
  SourceId: string;
  Company: string;
  Role: string;
  OfferUrl: string;
  Location: string
  SalaryMin: number;
  SalaryMax: number;
  Description: string;
  Score: number;
  PublishedAt: string;
}
