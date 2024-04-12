import { useQuery } from "@tanstack/react-query";
import { Article } from "../../articles/types";
import { Sources, QueryParams } from "../../main/types";
import { articlesAPI } from "../services/articles/articles.service";

export const useArticles = (source: Sources, queryParams: QueryParams) => {
  const query = useQuery<Article[]>({
    queryKey: ["articles", source, queryParams],
    queryFn: articlesAPI,
  });

  return { query, articles: query.data ?? [] };
};
