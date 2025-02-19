import { throttle } from 'lodash-es';
import { StopIcon } from '@heroicons/react/24/solid';
import { HTMLProps } from 'react';

function SendIcon(props: HTMLProps<SVGSVGElement>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
      {...props}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M5 10l7-7m0 0l7 7m-7-7v18"
      />
    </svg>
  );
}

export default function UserMessageInput({
  value,
  onChange,
  onSubmit,
  loading = false,
  onStop,
}: {
  value: string;
  onChange: (value: string) => void;
  onSubmit: () => void;
  loading?: boolean;
  onStop?: () => void;
}) {
  return (
    <>
      <div className="relative flex w-full rounded border border-[#D9D9D9] bg-[#FAFAFA] px-3 py-2">
        <textarea
          className="flex h-[6.375rem] w-full resize-none rounded-xl bg-transparent text-base placeholder-[#8C8C8C] outline-none lg:h-40"
          placeholder="Give answer to your agents."
          onChange={throttle((e) => onChange(e.target.value))}
          onKeyDown={
            loading
              ? undefined
              : (e) => {
                  if (e.key !== 'Enter' || !e.metaKey) {
                    return;
                  }
                  onSubmit();
                }
          }
          value={value}
        />
        <button
          onClick={() => {
            if (loading) {
              onStop?.();
            } else {
              onSubmit();
            }
          }}
          className="absolute bottom-4 right-4 flex size-7 items-center justify-center rounded-full bg-gray-400 text-white transition-colors hover:bg-gray-500 lg:size-12"
          disabled={value.trim() === ''}
        >
          {loading ? (
            <StopIcon className="size-4 lg:size-7" />
          ) : (
            <SendIcon className="size-4 lg:size-7" />
          )}
        </button>
      </div>
    </>
  );
}
