import MarkdownRenderer from '@/components/MarkdownRenderer';
import ProfileImage from '@/components/ProfileImage';
import { ChevronDown, ChevronRight } from 'lucide-react';
import {useMemo, useState} from 'react';
import classNames from 'classnames';

export default function BotChatBubble({
  botName,
  text = '',
  profileImageUrl,
  working = false,
  role,
}: {
  botName: string;
  text?: string;
  profileImageUrl: string;
  working?: boolean;
  role: string;
}) {
  const [collapse, setCollapse] = useState(false);
  const title = useMemo(() => {
    return `${botName} - ${role}`;
  }, [botName, role]);

  return (
    <div className="flex w-full flex-row items-start gap-4 bg-transparent">
      <div className="relative size-10 flex-shrink-0">
        <ProfileImage
          className="rounded-full"
          alt="bot_profile_image"
          src={profileImageUrl}
        />
      </div>
      <div className="flex w-full flex-col gap-y-1 pb-2">
        <div className="flex w-full items-center justify-between">
          <span className="flex text-sm font-bold">{title}</span>
          {!working && (
            <div className="flex flex-grow justify-end">
              <button
                className="flex items-center justify-center rounded p-0.5 hover:bg-gray-200"
                onClick={() => {
                  setCollapse((prev) => !prev);
                }}
              >
                {collapse ? (
                  <ChevronRight className="m-auto flex size-4" />
                ) : (
                  <ChevronDown className="m-auto flex size-4" />
                )}
              </button>
            </div>
          )}
        </div>
        {working && (
          <div className="flex w-full items-center gap-2">
            <div className="flex items-center space-x-1">
              <div className="flex space-x-1">
                {[0, 1, 2].map((index) => (
                  <div
                    key={index}
                    className="h-2 w-2 animate-bounce rounded-full bg-gray-400"
                    style={{
                      animationDelay: `${index * 0.2}s`,
                      animationDuration: '0.8s',
                    }}
                  />
                ))}
              </div>
            </div>
          </div>
        )}
        {!working && (
          <MarkdownRenderer
            className={classNames('prose prose-sm prose-a:text-blue-600', {
              'line-clamp-1': collapse,
            })}
            content={text}
          />
        )}
      </div>
    </div>
  );
}
