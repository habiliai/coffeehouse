import { Children, isValidElement, PropsWithChildren } from 'react';
import ProfileImage from './ProfileImage';
import classNames from 'classnames';
import { AgentWork } from '@/proto/habapi';

function AgentProfile({
  className,
  imageClassName,
  name,
  imageUrl,
  children,
  status,
}: PropsWithChildren<{
  className?: string;
  imageClassName?: string;
  name: string;
  imageUrl: string;
  status?: AgentWork.Status;
}>) {
  const label = Children.toArray(children).find((child) => {
    return isValidElement(child) && child.type === AgentProfile.Label;
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
          {status && <div className={classNames("absolute top-0 text-red-600 font-extrabold", {
            '-right-4': status === AgentWork.Status.IDLE,
            'right-0': status === AgentWork.Status.WAITING,
            '-right-1.5': status === AgentWork.Status.WORKING,
          })}>{
            status === AgentWork.Status.WORKING ? "W" :
            status === AgentWork.Status.IDLE ? "Zzz" :
            status === AgentWork.Status.WAITING ? "!" : ""
          }</div>}
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
