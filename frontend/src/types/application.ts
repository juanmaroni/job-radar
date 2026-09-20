export interface Application {
  Id: number;
  Company: string;
  Role: string;
  OfferUrl: string;
  Status: string;
  MatchScore: number | null;
  CreatedAt: string;
  UpdatedAt: string;
}