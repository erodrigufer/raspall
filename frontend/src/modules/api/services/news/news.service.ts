import { Article } from "../../../articles/types";
import { generateGetMethod } from "../utils";

export const newsURL = "/v1/news";

export const newsFunction = generateGetMethod<Article[]>(newsURL);
