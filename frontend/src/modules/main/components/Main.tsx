import { useQuery } from "@tanstack/react-query";
import Articles from "../../articles/components/Articles";
import { Article } from "../../articles/types";
import { newsFunction } from "../../api/services/news/news.service";

const Main: React.FC = () => {
  const query = useQuery<Article[]>({
    queryKey: ["news", "zeit"],
    queryFn: newsFunction,
    refetchOnWindowFocus: false,
  });

  const renderedArticles = query?.data;

  return renderedArticles?.length !== 0 && renderedArticles ? (
    <Articles articles={renderedArticles} />
  ) : (
    <p>No new articles</p>
  );
};

export default Main;
