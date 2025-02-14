import { PlaceholderValue } from 'next/dist/shared/lib/get-img-props';
import { toBase64 } from '@/lib/base64';
import Image from 'next/image';

function shimmer(w: number | string, h: number | string) {
  return `<svg width="${w}" height="${h}" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <defs>
    <linearGradient id="g">
      <stop stop-color="#333" offset="20%" />
      <stop stop-color="#222" offset="50%" />
      <stop stop-color="#333" offset="70%" />
    </linearGradient>
  </defs>
  <rect width="${w}" height="${h}" fill="#333" />
  <rect id="r" width="${w}" height="${h}" fill="url(#g)" />
  <animate xlink:href="#r" attributeName="x" from="-${w}" to="${w}" dur="1s" repeatCount="indefinite"  />
</svg>`;
}

const placeholder: PlaceholderValue = `data:image/svg+xml;base64,${toBase64(shimmer('100%', '100%'))}`;

const ProfileImage = ({
  className,
  alt,
  src,
}: {
  className?: string;
  alt: string;
  src: string;
}) => {
  src = src ?? placeholder;
  if (src === '') src = placeholder;

  return (
    <Image
      fill
      className={className}
      alt={alt}
      src={src}
      width={0}
      height={0}
      sizes="100vw"
      placeholder={placeholder}
    />
  );
};

export default ProfileImage;
