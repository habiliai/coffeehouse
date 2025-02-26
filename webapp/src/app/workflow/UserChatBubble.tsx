import { Fragment } from 'react';

export default function UserChatBubble({
  text,
  id,
}: {
  id: string;
  text: string;
}) {
  return (
    <div className="flex flex-row-reverse">
      <div className="flex max-w-[75%] flex-wrap items-center justify-start gap-x-1 gap-y-0.5 rounded-md bg-[#f3f3f3] p-4 text-sm text-gray-900">
        {text.split('\n').map((line, i) => (
          <Fragment key={`user-messsage-${id}-line-${i}`}>
            {line.split(' ').map((word, j) => {
              const key = `user-message-${id}-line-${i}-word-${j}`;
              return word.startsWith('@') ? (
                <span
                  key={key}
                  className="font-semibold text-blue-500 underline"
                >
                  {word}
                </span>
              ) : (
                <span key={key}>{word}</span>
              );
            })}
            <br />
          </Fragment>
        ))}
      </div>
    </div>
  );
}
