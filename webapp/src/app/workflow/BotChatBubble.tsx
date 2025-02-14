import MarkdownRenderer from '@/components/MarkdownRenderer';
import ProfileImage from '@/components/ProfileImage';
import { ChevronRight } from 'lucide-react';

export default function BotChatBubble({
  botName,
  text,
  profileImageUrl,
}: {
  botName: string;
  text: string;
  profileImageUrl: string;
}) {
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
        <div className="flex w-full justify-between items-center">
          <span className="flex text-sm font-bold">{botName}</span>
          <div className="flex flex-grow justify-end">
            <ChevronRight className="size-4" />
          </div>
        </div>
        <MarkdownRenderer className="prose prose-sm" content={text} />
      </div>
    </div>
  );
}
