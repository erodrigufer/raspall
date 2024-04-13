import { List, ListItem, ListItemText, Paper } from "@mui/material";
import { Article } from "../types";

type Props = {
  articles: Article[];
};

const Articles: React.FC<Props> = ({ articles }) => {
  return (
    <>
      <Paper elevation={1}>
        <List>
          {articles.map((article) => (
            <ListItem component="a" href={article.URL} key={article.Title}>
              <ListItemText primary={article.Title} />
            </ListItem>
          ))}
        </List>
      </Paper>
    </>
  );
};

export default Articles;
