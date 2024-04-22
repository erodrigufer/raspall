import { QueryParams } from "./types";

export function NewQueryParams(
  limit: number,
  removePaywall: boolean = true,
): QueryParams {
  return { limit: limit, removePaywall: removePaywall };
}
