import './Banner.css';

interface BannerProps {
  src: string;
  alt: string;
}

export default function Banner({ src, alt }: BannerProps) {
  return (
    <div className="banner">
      <img className="banner__image" src={src} alt={alt} fetchPriority="high" />
    </div>
  );
}
