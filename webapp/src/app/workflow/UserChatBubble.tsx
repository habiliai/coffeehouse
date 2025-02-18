import { Fragment } from 'react';

export default function UserChatBubble({
  text,
}: {
  text: string;
  mentions: string[];
}) {
  return (
    <div className="flex flex-row-reverse">
      <div className="flex justify-start items-center max-w-[75%] rounded-md bg-[#f3f3f3] p-4 text-sm text-gray-900 gap-y-0.5 gap-x-1 flex-wrap">
        {text.split('\n').map((line) => (
          <>
            {line.split(' ').map((word) => (
              <>
                {word.startsWith('@') ? (
                  <span className="font-semibold text-blue-500 underline">
                    {word}
                  </span>
                ) : (
                  <span>{word}</span>
                )}
              </>
            ))}
            <br />
          </>
        ))}
      </div>
    </div>
  );
}
