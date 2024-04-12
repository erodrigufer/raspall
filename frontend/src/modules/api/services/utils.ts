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
    const [, source, params] = context.queryKey;
    const queryParams = new URLSearchParams(
      removeUndefinedValues(params as Record<string, string | undefined>),
    ).toString();
    const requestURL = `${url}/${source}?${queryParams}`;
    console.info(`Fetching ${requestURL}`);
    const response = await fetch(new URL(requestURL, "http://localhost"), {
      method: "GET",
      credentials: "same-origin",
    });

    checkStatusCode(response);

    return await response.json();
  };
};

const checkStatusCode = (response: Response) => {
  if (!response.ok) {
    throw new Error(
      `Failed to fetch ${response.url}: ${response.status} ${response.statusText}`,
    );
  }
};

// Returns an object with no undefined values and all keys are of type 'string',
// i.e. Record<string, string>.
// Record<K, T> is a utility type that constructs an object type whose keys are of type K
// and values are of type T.
// It's a generic type that provides a way to declare the shape of an object when
// the exact property names are not important. In this case the input object will have keys
// of type 'string' and its values might be 'string' or 'undefined'.
const removeUndefinedValues = (
  params: Record<string, string | undefined>,
): Record<string, string> => {
  // Object.fromEntries takes an array of key-value pairs and constructs an object
  // from them.
  return Object.fromEntries(
    // Object.entries creates an array of `[key, value]` pairs.
    // Using Object.entries without spreading params first touches
    // the original object and triggers a re-render in react-query.
    // Spreading params makes a shallow-copy of params.
    Object.entries({ ...params }).filter(notUndefinedEntry),
  );
};

// notUndefinedEntry is a type guard that returns true if
// the value in a [key value] pair is not undefined.
const notUndefinedEntry = <T>(
  entry: [string, T | undefined],
): entry is [string, T] => {
  return entry[1] !== undefined;
};
