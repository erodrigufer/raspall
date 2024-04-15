import { useQuery } from "@tanstack/react-query";
import Articles from "../../articles/components/Articles";
import { Article } from "../../articles/types";
import { newsFunction } from "../../api/services/news/news.service";
import { Box } from "@mui/material";
import SourceBadge from "../../sources/components/SourceBadge";

const Main: React.FC = () => {
  const query = useArticles("zeit");
  const renderedArticles = query?.data;

  return renderedArticles?.length !== 0 && renderedArticles ? (
    <Box sx={{ width: "100%", maxWidth: "40vw" }}>
      <SourceBadge />
      <Articles articles={renderedArticles} />
    </Box>
  ) : (
    <p>No new articles</p>
  );
};

export default Main;

export type QueryParams = {
  limit: number;
  removePaywall: boolean;
};

export const useArticles = (source: string) => {
  const params: QueryParams = {
    limit: 5,
    removePaywall: true,
  };
  const query = useQuery<Article[]>({
    queryKey: ["articles", source, params],
    queryFn: newsFunction,
    refetchOnWindowFocus: false,
  });

  return { ...query, articles: query.data ?? [] };
};

//   const query = useQuery({
//     queryKey: ['statistics', { category, mode, interval, from, until, objectId }],
//     queryFn: (context) => api.statistics.find(context),
//   });
// };
