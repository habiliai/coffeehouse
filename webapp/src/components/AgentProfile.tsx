import { Children, isValidElement, PropsWithChildren } from 'react';
import { PlaceholderValue } from 'next/dist/shared/lib/get-img-props';
import { toBase64 } from '@/lib/base64';
import Image from 'next/image';

function AgentProfile({
  name,
  imageUrl,
  children,
}: PropsWithChildren<{
  name: string;
  imageUrl: string;
}>) {
  const label = Children.toArray(children).find((child) => {
    return isValidElement(child) && child.type === AgentProfile.Label;
  });

  const icon = Children.toArray(children).find((child) => {
    return isValidElement(child) && child.type === AgentProfile.Icon;
  });

  return (
    <div className="flex w-32 flex-col items-center gap-y-2 p-2">
      <div className="relative size-16 overflow-hidden rounded-full">
        <AgentImage alt={`agent-${name}`} src={imageUrl} />
        {icon && <div className="absolute bottom-0 right-0">{icon}</div>}
      </div>
      {label}
    </div>
  );
}

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

const AgentImage = ({ alt, src }: { alt: string; src: string }) => {
  src = src ?? placeholder;
  if (src === '') src = placeholder;

  return (
    <Image
      alt={alt}
      src={src}
      width={64}
      height={64}
      sizes="100vw"
      placeholder={placeholder}
    />
  );
};

const AgentLabel = ({
  className,
  children,
}: PropsWithChildren<{ className?: string }>) => {
  return <label className={className}>{children}</label>;
};
AgentProfile.Label = AgentLabel;

const AgentIcon = ({
  className,
  children,
}: PropsWithChildren<{ className?: string }>) => {
  return <div className={className}>{children}</div>;
};
AgentProfile.Icon = AgentIcon;

export default AgentProfile;
