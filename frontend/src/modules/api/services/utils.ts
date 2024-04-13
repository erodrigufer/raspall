import { QueryFunctionContext } from "@tanstack/react-query";

/**
 * Generate a function that fetches a single item from the API.
 * The fetch call uses GET.
 *
 * @param url The URL to fetch from.
 * @returns A function that can be used to fetch an item with the given id.
 */
export const generateGetMethod = <T>(url: string) => {
  return async (context: QueryFunctionContext): Promise<T> => {
    const [, id] = context.queryKey;
    console.info(`Fetching ${url}/${id} (get)`);
    const response = await fetch(
      new URL(`${url}/${id}?limit=10`, "http://localhost"),
      {
        method: "GET",
        credentials: "omit",
      },
    );

    // checkStatusCode(response, context.meta?.ignorePermissionDenied);

    return await response.json();
  };
};
