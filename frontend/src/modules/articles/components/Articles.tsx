import { List, ListItem, ListItemButton, ListItemText } from "@mui/material";
import { Article } from "../types";

type Props = {
  articles: Article[];
};

const Articles: React.FC<Props> = ({ articles }) => {
  return (
    <>
      <List>
        {articles.map((article) => (
          <ListItem disablePadding>
            <ListItemButton component="a" href={article.URL}>
              <ListItemText primary={article.Title} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </>
  );
};

export default Articles;
