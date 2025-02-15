'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useAddMessage, useGetThread } from '@/actions/thread';
import AgentProfile from '@/components/AgentProfile';
import {
  ChevronDown,
  ChevronLeft,
  ChevronUp,
  CircleCheck,
  Ellipsis,
} from 'lucide-react';
import { useCallback, useState } from 'react';
import UserChatBubble from './UserChatBubble';
import BotChatBubble from './BotChatBubble';
import UserMessageInput from './UserMessageInput';
import { Button } from '@/components/ui/button';
import LoadingSpinner from '@/components/LoadingSpinner';
import MarkdownRenderer from '@/components/MarkdownRenderer';
import classNames from 'classnames';
import ProfileImage from '@/components/ProfileImage';

export default function Page() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const threadId = searchParams.get('thread_id') ?? '';

  const [isRunning, setRunning] = useState<boolean>(false);
  const [userMessage, setUserMessage] = useState('');
  const [selectedStep, setSelectedStep] = useState<number>(0);
  const [resultCollapsed, setResultCollapsed] = useState<boolean>(true);

  const { data } = useGetThread({ threadId: threadId });

  const { mutate: addMessage } = useAddMessage({
    onSuccess: () => {
      setRunning(false);
    },
    onError: () => {
      setRunning(false);
    },
    onMutate: (message) => {
      data?.thread.messagesList.push({
        id: Date.now().toString(),
        role: 1,
        text: message,
        mentionsList: [],
      });
    },
  });

  const handleOnSubmitUserMessage = useCallback(() => {
    addMessage({ message: userMessage });
    setUserMessage('');
  }, [addMessage, userMessage, setUserMessage]);

  const handleCopyAsMarkdown = useCallback(() => {
    // TODO: Implement copy as markdown logic
  }, []);

  return (
    <div className="flex w-full flex-row gap-[0.875rem] bg-[#F7F7F7]">
      <div className="hidden w-full max-w-[8.125rem] flex-col items-center gap-2 border border-[#E5E7EB] bg-white py-[2.375rem] shadow-[0_0_40px_3px_#AEAEAE40] lg:flex">
        <button className="mb-14" onClick={() => router.back()}>
          <ChevronLeft />
        </button>
        <div className="mb-7">
          <span className="text-sm font-bold">Agents</span>
        </div>
        {!data && <LoadingSpinner className="flex h-12 w-12" />}
        {data?.agents.map((agent) => (
          <AgentProfile
            key={`workflow-agent-${agent.id}`}
            name={agent.name}
            imageUrl={agent.iconUrl}
            imageClassName="size-10"
          >
            <AgentProfile.Label className="text-sm font-light">
              {agent.name}
            </AgentProfile.Label>
          </AgentProfile>
        ))}
      </div>

      <div className="flex h-full w-full max-w-[30.25rem] flex-col border border-[#E5E7EB] bg-white shadow-[0_0_40px_3px_#AEAEAE40]">
        <div className="relative flex justify-center px-14 py-6">
          <button
            className="absolute left-2 flex lg:hidden"
            onClick={() => router.back()}
          >
            <ChevronLeft className="text-black" />
          </button>
          <span className="line-clamp-1 text-sm font-normal">
            {data?.thread.title}
          </span>
        </div>

        <div className="flex h-full flex-col gap-4 overflow-y-auto border-t border-[#E2E8F0] px-6 py-9">
          {!data && <LoadingSpinner className="m-auto flex h-12 w-12" />}
          {data?.thread.messagesList.map(({ id, role, text, agent }) => {
            switch (role) {
              case 1:
                return <UserChatBubble key={`thread-user-${id}`} text={text} />;
              case 2:
                return (
                  <BotChatBubble
                    key={`thread-bot-${id}`}
                    botName={agent?.name || ''}
                    text={text}
                    profileImageUrl={agent?.iconUrl || ''}
                  />
                );
            }
          })}
        </div>
        <div className="relative flex w-full flex-grow items-end px-6 pb-[1.875rem]">
          <UserMessageInput
            loading={isRunning || !data}
            value={userMessage}
            onChange={(v) => setUserMessage(v)}
            onSubmit={handleOnSubmitUserMessage}
          />
        </div>
      </div>

      <div className="hidden w-full max-w-[calc(100%-8.125rem-30.25rem)] flex-col items-center gap-10 pb-[0.875rem] pr-[0.875rem] pt-16 lg:flex">
        <div className="flex w-full max-w-[50.5625rem] flex-col items-center justify-center gap-y-4">
          {!data && <LoadingSpinner className="m-auto flex h-12 w-12" />}
          <div className="flex w-full p-[0.3125rem]">
            {data?.thread.stepsList.map((_, index) => (
              <WorkflowStepButton
                key={`workflow-step-${index}`}
                currentStep={data.thread.currentStep}
                index={index}
                selected={selectedStep === index}
                onClick={() => setSelectedStep(index)}
              />
            ))}
          </div>
          <div className="flex w-full flex-col gap-1">
            {data?.thread.stepsList[selectedStep].tasksList.map(
              (task, index) => (
                <div
                  key={`workflow-task-${index}`}
                  className="flex w-full items-center justify-between rounded-[0.625rem] border border-[#E2E8F0] bg-white py-[0.3125rem] pl-[1.875rem] pr-[0.3125rem]"
                >
                  <div className="flex gap-[1.875rem]">
                    {task.status === 'done' ? (
                      <CircleCheck className="text-black" />
                    ) : (
                      <Ellipsis className="text-black" />
                    )}
                    <span className="line-clamp-1 text-sm/[1.375rem] font-normal">
                      {task.title}
                    </span>
                  </div>
                  <div className="flex">
                    {task.requiredAgents.map((agent) => (
                      <div
                        key={`workflow-task-agent-${agent.id}`}
                        className="relative size-10 overflow-hidden rounded-full"
                      >
                        <ProfileImage alt={agent.name} src={agent.iconUrl} />
                      </div>
                    ))}
                  </div>
                </div>
              ),
            )}
          </div>
        </div>
        <div className="flex w-full flex-1 flex-col justify-between rounded-[1.25rem] border border-[#E5E7EB] bg-white p-4 shadow-[0_0_40px_3px_#AEAEAE40]">
          <div className="flex w-full flex-col gap-[0.875rem] px-[1.875rem] py-[1.625rem]">
            <span className="text-[2rem]/[3.375rem] font-bold">Result</span>
            <div
              className={classNames(
                'relative flex flex-1 flex-col overflow-hidden',
                {
                  'max-h-[calc(20.25rem)] overflow-auto': resultCollapsed,
                },
              )}
            >
              <div
                className={classNames(
                  'absolute z-50 flex h-full w-full items-end justify-center',
                  {
                    'bg-gradient-to-b from-transparent to-[#FFFFFF]':
                      resultCollapsed,
                  },
                )}
              >
                {resultCollapsed ? (
                  <ChevronDown
                    className="size-10 text-[#878787]"
                    onClick={() => setResultCollapsed(false)}
                  />
                ) : (
                  <ChevronUp
                    className="size-10 text-[#878787]"
                    onClick={() => setResultCollapsed(true)}
                  />
                )}
              </div>
              <MarkdownRenderer content={data?.thread.resultContent ?? ''} />
            </div>
          </div>

          <div className="flex w-full">
            <Button
              className="w-full"
              disabled={data?.thread.status !== 'done'}
              onClick={handleCopyAsMarkdown}
            >
              Copy as Markdown
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}

function WorkflowStepButton({
  currentStep,
  selected,
  index,
  onClick,
}: {
  currentStep: number;
  selected: boolean;
  index: number;
  onClick?: () => void;
}) {
  return (
    <button
      className={classNames('rounded-[0.1875rem] px-[0.75rem] py-[0.375rem]', {
        'bg-white font-bold text-[#0F172A]': selected,
        'bg-transparent font-medium text-[#334155]': !selected,
        'cursor-not-allowed': index > currentStep - 1,
      })}
      onClick={() => {
        if (index > currentStep - 1) return;
        onClick?.();
      }}
    >
      Step {index + 1}
    </button>
  );
}
