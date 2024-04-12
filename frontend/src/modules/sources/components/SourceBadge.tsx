import { Chip } from "@mui/material";

type Props = {
  label: string;
  onClick: () => void;
  selected?: boolean;
};

const SourceBadge: React.FC<Props> = ({ label, onClick, selected }) => {
  return (
    <Chip
      onClick={onClick}
      label={label}
      color={selected ? "primary" : undefined}
    />
  );
};

export default SourceBadge;
