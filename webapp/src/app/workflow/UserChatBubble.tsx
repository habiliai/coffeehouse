import { Fragment } from 'react';

export default function UserChatBubble({
  text,
  mentions,
}: {
  text: string;
  mentions: string[];
}) {
  return (
    <div className="flex flex-row-reverse">
      <div className="flex max-w-[75%] rounded-md bg-[#f3f3f3] p-4 text-sm text-gray-900">
        {mentions.map((mention) => (
          <Fragment key={`${mention}`}>
            <span className="font-semibold text-blue-500 underline">
              @{mention}
            </span>
            &nbsp;
          </Fragment>
        ))}
        {text}
      </div>
    </div>
  );
}
