import Articles from "../../articles/components/Articles";
import {
  Box,
  Card,
  CircularProgress,
  Container,
  Divider,
  Stack,
  Typography,
} from "@mui/material";
import SourceBadge from "../../sources/components/SourceBadge";
import { useState } from "react";
import { QueryParams, Sources } from "../types";
import { NewQueryParams } from "../utils";
import { useArticles } from "../../api/hooks/useArticles";

const defaultSource = "nació";
const defaultQueryParams = {
  limit: 10,
  removePaywall: true,
};

const Main: React.FC = () => {
  const [source, setSource] = useState<Sources>(defaultSource);
  const [queryParams, setQueryParams] =
    useState<QueryParams>(defaultQueryParams);

  const onClickSourceBadge =
    (source: Sources, queryParams: QueryParams) => () => {
      setSource(source);
      setQueryParams(queryParams);
    };

  const { query, articles } = useArticles(source, queryParams);
  const renderedArticles = articles;

  return (
    <Container maxWidth="lg">
      <Card variant="outlined" sx={{ width: "100%" }}>
        <Stack
          direction="row"
          spacing={1}
          sx={{ marginBottom: "1em", p: "1em" }}
        >
          <SourceBadge
            onClick={onClickSourceBadge("nació", NewQueryParams(10))}
            label={"Nació"}
            selected={source === "nació" ? true : undefined}
          />
          <SourceBadge
            onClick={onClickSourceBadge("hn", NewQueryParams(30))}
            label={"Hacker News"}
            selected={source === "hn" ? true : undefined}
          />
          <SourceBadge
            onClick={onClickSourceBadge("zeit", NewQueryParams(10))}
            label={"Zeit"}
            selected={source === "zeit" ? true : undefined}
          />
        </Stack>
        <Divider />

        <Box sx={{ p: "1em" }}>
          {renderedArticles?.length !== 0 ? (
            <Articles articles={renderedArticles} />
          ) : query.isPending ? (
            <Stack direction="row" spacing={1}>
              <CircularProgress />
              <Typography> Loading... </Typography>
            </Stack>
          ) : (
            <Typography>No new articles</Typography>
          )}
        </Box>
      </Card>
    </Container>
  );
};

export default Main;
