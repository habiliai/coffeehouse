import { Children, isValidElement, PropsWithChildren } from 'react';
import ProfileImage from './ProfileImage';
import classNames from 'classnames';

function AgentProfile({
  className,
  imageClassName,
  name,
  imageUrl,
  children,
}: PropsWithChildren<{
  className?: string;
  imageClassName?: string;
  name: string;
  imageUrl: string;
}>) {
  const label = Children.toArray(children).find((child) => {
    return isValidElement(child) && child.type === AgentProfile.Label;
  });

  const icon = Children.toArray(children).find((child) => {
    return isValidElement(child) && child.type === AgentProfile.Icon;
  });

  const tooltip = Children.toArray(children).find((child) => {
    return isValidElement(child) && child.type === AgentProfile.Tooltip;
  });

  return (
    <div className={className}>
      <div className="group relative flex flex-col items-center gap-y-2">
        <div className={classNames('relative', imageClassName)}>
          <ProfileImage
            className="rounded-full"
            alt={`agent-${name}`}
            src={imageUrl}
          />
          {icon && <div className="absolute bottom-0 right-0">{icon}</div>}
        </div>
        {label}

        {tooltip && (
          <div className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 flex w-52 -translate-x-1/2 rounded-lg bg-gray-700 p-2 opacity-0 ring-1 ring-inset ring-gray-400 transition-opacity group-hover:opacity-100">
            {tooltip}
          </div>
        )}
      </div>
    </div>
  );
}

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

const AgentTooltip = ({
  className,
  children,
}: PropsWithChildren<{ className?: string }>) => {
  return <div className={className}>{children}</div>;
};
AgentProfile.Tooltip = AgentTooltip;

export default AgentProfile;
