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
      <div className="flex flex-col gap-y-1 pb-2">
        <div className="flex w-full justify-between">
          <span className="text-sm font-bold">{botName}</span>
          <ChevronRight className="size-6" />
        </div>
        <MarkdownRenderer className="prose prose-sm" content={text} />
      </div>
    </div>
  );
}
