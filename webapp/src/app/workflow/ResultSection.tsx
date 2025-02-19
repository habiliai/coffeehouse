import MarkdownRenderer from '@/components/MarkdownRenderer';
import { Button } from '@/components/ui/button';
import classNames from 'classnames';
import { ChevronDown, ChevronUp } from 'lucide-react';

export default function ResultSection({
  markdownContent,
  resultCollapsed,
  setResultCollapsed,
  onClickCopyButton,
}: {
  markdownContent: string;
  resultCollapsed: boolean;
  setResultCollapsed: (collapsed: boolean) => void;
  onClickCopyButton: () => void;
}) {
  return (
    <div className="flex h-full w-full flex-1 flex-col justify-between">
      <div className="relative flex flex-1 flex-col items-center gap-[0.875rem] p-0 lg:overflow-hidden lg:px-[1.875rem] lg:py-[1.625rem]">
        {markdownContent && resultCollapsed && (
          <div
            className={classNames(
              'absolute z-50 hidden h-full w-full items-end justify-center pb-5 lg:flex',
              'bg-gradient-to-b from-transparent to-[#FFFFFF]',
            )}
          >
            <ChevronDown
              className="size-14 cursor-pointer text-[#878787]"
              onClick={() => setResultCollapsed(false)}
            />
          </div>
        )}

        <span className="w-full text-2xl/[3.375rem] font-bold text-black">
          Result
        </span>
        <div className="relative w-full flex-1">
          <MarkdownRenderer
            className="prose prose-sm"
            content={markdownContent}
          />
        </div>
      </div>
      {!resultCollapsed && (
        <div className="flex w-full justify-center">
          <ChevronUp
            className="size-14 cursor-pointer text-[#878787]"
            onClick={() => setResultCollapsed(true)}
          />
        </div>
      )}
      <div className="flex w-full">
        {markdownContent && (
          <Button className="w-full" onClick={onClickCopyButton}>
            Copy as Markdown
          </Button>
        )}
      </div>
    </div>
  );
}
