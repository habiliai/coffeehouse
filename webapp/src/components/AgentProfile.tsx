'use client';

import { Children, isValidElement, PropsWithChildren } from 'react';
import ProfileImage from './ProfileImage';
import classNames from 'classnames';
import { AgentWork } from '@/proto/habapi';

function AgentProfile({
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
    <div className="group relative flex flex-col items-center gap-y-2 justify-start">
      <div className="relative flex">
        <div
          className={classNames(
            'relative min-h-6 min-w-6 rounded-full ' + imageClassName,
            {
              'border border-[#CBD5E1]': status === AgentWork.Status.IDLE,
              'border border-[#4AD15B]': status === AgentWork.Status.WAITING,
            },
          )}
        >
          <ProfileImage
            className={classNames({
              'opacity-50': status === AgentWork.Status.IDLE,
            })}
            alt={`agent-${name}`}
            src={imageUrl}
          />
        </div>
        <div
          className={classNames(
            'absolute -right-[0.125rem] bottom-0 rounded-full border bg-white p-1',
            {
              hidden: !status || status === AgentWork.Status.WORKING,
              'border-[#CBD5E1]': status === AgentWork.Status.IDLE,
              'border-[#4AD15B]': status === AgentWork.Status.WAITING,
            },
          )}
        >
          <StatusIcon status={status} />
        </div>
      </div>
      <span>{label}</span>

      {tooltip && (
        <div className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 flex w-52 -translate-x-1/2 rounded-lg bg-gray-700 p-2 opacity-0 ring-1 ring-inset ring-gray-400 transition-opacity group-hover:opacity-100">
          {tooltip}
        </div>
      )}
    </div>
  );
}

function StatusIcon({ status }: { status?: AgentWork.Status }) {
  switch (status) {
    case AgentWork.Status.IDLE:
      return <div className="text-[7px]">💤</div>;
    case AgentWork.Status.WAITING:
      return <div className="text-[7px]">❗️</div>;
    default:
      return null;
  }
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
