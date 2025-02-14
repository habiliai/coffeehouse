import { throttle } from 'lodash-es';
import { StopIcon } from '@heroicons/react/24/solid';
import { SendIcon } from '@/components/Icons';

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
          className="flex h-40 w-full resize-none rounded-xl bg-transparent text-sm/[1.5rem] placeholder-[#8C8C8C] outline-none"
          placeholder="Give answer to your agents."
          value={value}
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
        />
        <button
          onClick={() => {
            if (loading) {
              onStop?.();
            } else {
              onSubmit();
            }
          }}
          className="absolute bottom-4 right-4 flex size-12 items-center justify-center rounded-full bg-gray-400 text-white transition-colors hover:bg-gray-500"
        >
          {loading ? (
            <StopIcon className="size-7" />
          ) : (
            <SendIcon className="size-7" />
          )}
        </button>
      </div>
    </>
  );
}
