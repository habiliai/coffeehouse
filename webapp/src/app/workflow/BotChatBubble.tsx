import MarkdownRenderer from '@/components/MarkdownRenderer';
import ProfileImage from '@/components/ProfileImage';
import classNames from 'classnames';
import { ChevronDown, ChevronRight } from 'lucide-react';
import { useState } from 'react';

export default function BotChatBubble({
  botName,
  text,
  profileImageUrl,
}: {
  botName: string;
  text: string;
  profileImageUrl: string;
}) {
  const [isCollapsed, setCollapsed] = useState<boolean>(true);
  return (
    <button
      className="flex w-full flex-row items-start gap-4 bg-transparent text-left"
      onClick={() => setCollapsed(!isCollapsed)}
    >
      <div className="relative size-10 flex-shrink-0">
        <ProfileImage
          className="rounded-full"
          alt="bot_profile_image"
          src={profileImageUrl}
        />
      </div>
      <div className="flex w-full flex-col gap-y-1 pb-2">
        <div className="flex w-full items-center justify-between">
          <span className="flex text-sm font-bold">{botName}</span>
          <div className="flex flex-grow justify-end">
            {isCollapsed ? (
              <ChevronRight className="size-4" />
            ) : (
              <ChevronDown className="size-4" />
            )}
          </div>
        </div>
        <MarkdownRenderer
          className={classNames('prose prose-sm', {
            'line-clamp-1': isCollapsed,
          })}
          content={text}
        />
      </div>
    </button>
  );
}
