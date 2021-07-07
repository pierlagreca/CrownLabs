import { Property } from 'csstype';
import { FC } from 'react';
import './CrownLoader.css';

export interface ICrownLoaderProps {
  duration?: number;
  color?: Property.Stroke;
  size?: number | string;
}

const Loader: FC<ICrownLoaderProps> = ({ ...props }) => {
  const { duration = 4, color = '#FFF', size } = props;
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      height={size}
      width={size}
      viewBox="0 0 512 512"
    >
      <path
        className="crown"
        style={{
          stroke: color,
          animationDuration: `${duration}s`,
        }}
        d="M170 292h172l42-128-80 48-48-112-48 112-80-48 42 128z"
      />
      <path
        style={{
          stroke: color,
          animationDuration: `${duration * 0.7}s`,
        }}
        className="dash"
        d="M170 320h172"
      />
    </svg>
  );
};

export default Loader;
