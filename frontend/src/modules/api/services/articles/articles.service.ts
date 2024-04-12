import { Article } from "../../../articles/types";
import { generateGetMethod } from "../utils";

export const articlesURL = "/v1/articles";

export const articlesAPI = generateGetMethod<Article[]>(articlesURL);
