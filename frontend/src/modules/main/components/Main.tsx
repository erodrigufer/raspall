import { useQuery } from "@tanstack/react-query";
import Articles from "../../articles/components/Articles";
import { Article } from "../../articles/types";
import { newsFunction } from "../../api/services/news/news.service";
import { Box } from "@mui/material";
import SourceBadge from "../../sources/components/SourceBadge";
import { useState } from "react";

type Sources = "nació" | "zeit" | "hn";

const Main: React.FC = () => {
  const [source, setSource] = useState<Sources>("nació");
  const [queryParams, setQueryParams] = useState<QueryParams>({
    limit: 10,
    removePaywall: true,
  });

  const onClickSourceBadge =
    (source: Sources, queryParams: QueryParams) => () => {
      setSource(source);
      setQueryParams(queryParams);
    };

  const query = useArticles(source, queryParams);
  const renderedArticles = query?.data;

  return (
    <Box sx={{ width: "100%", maxWidth: "40vw" }}>
      <SourceBadge
        onClick={onClickSourceBadge("nació", NewQueryParams(10))}
        label={"Nació"}
      />
      <SourceBadge
        onClick={onClickSourceBadge("hn", NewQueryParams(30))}
        label={"Hacker News"}
      />
      <SourceBadge
        onClick={onClickSourceBadge("zeit", NewQueryParams(10))}
        label={"Zeit"}
      />

      {renderedArticles?.length !== 0 && renderedArticles ? (
        <Articles articles={renderedArticles} />
      ) : query.isPending ? (
        <p> Loading... </p>
      ) : (
        <p>No new articles</p>
      )}
    </Box>
  );
};

export default Main;

export type QueryParams = {
  limit: number;
  removePaywall: boolean;
};

function NewQueryParams(
  limit: number,
  removePaywall: boolean = true,
): QueryParams {
  return { limit: limit, removePaywall: removePaywall };
}

export const useArticles = (source: Sources, queryParams: QueryParams) => {
  const query = useQuery<Article[]>({
    queryKey: ["articles", source, queryParams],
    queryFn: newsFunction,
  });

  return { ...query, articles: query.data ?? [] };
};
