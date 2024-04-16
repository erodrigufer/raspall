import { Chip } from "@mui/material";

type Props = {
  label: string;
  onClick: () => void;
};

const SourceBadge: React.FC<Props> = ({ label, onClick }) => {
  return <Chip onClick={onClick} size="small" label={label} />;
};

export default SourceBadge;
