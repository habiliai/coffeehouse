'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import {
  useAddMessage,
  useGetStatus,
  useGetThread,
  useGetMissionStepStatus,
} from './actions';
import AgentProfile from '@/components/AgentProfile';
import { ChevronLeft, ChevronDown } from 'lucide-react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import UserMessageInput from './UserMessageInput';
import LoadingSpinner from '@/components/LoadingSpinner';
import { Agent, AgentWork } from '@/proto/habapi';
import classNames from 'classnames';
import ResultSection from './ResultSection';
import WorkflowSection from './WorkflowSection';
import MobileNavbar from './MobileNavbar';
import ChatSection from './ChatSection';

export default function Page() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const threadId = useMemo(() => {
    const threadId = searchParams.get('thread_id');
    return threadId ? parseInt(threadId) : null;
  }, [searchParams]);

  const {
    data: { thread, mission, lastMessageId },
  } = useGetThread({ threadId });

  const {
    data: { agentWorks },
  } = useGetStatus({
    threadId,
    lastMessageId,
  });

  const [isRunning, setRunning] = useState<boolean>(false);
  const [userMessage, setUserMessage] = useState('');
  const [nowDisplayedStep, setNowDisplayedStep] = useState<number | null>(null);
  const [resultCollapsed, setResultCollapsed] = useState<boolean>(true);
  const [mobileView, setMobileView] = useState<'Chat' | 'Workflow' | 'Outcome'>(
    'Chat',
  );

  const chatContainerRef = useRef<HTMLDivElement>(null);
  const observerRef = useRef<HTMLDivElement>(null);
  const [isAtBottom, setIsAtBottom] = useState(true);

  const {
    data: { initialLoading: isActionWorksLoading, actionWorks },
  } = useGetMissionStepStatus({
    threadId,
    selectedSeqNo: nowDisplayedStep,
  });

  const { mutate: addMessage } = useAddMessage({
    threadId,
    onSuccess: () => {
      setRunning(false);
    },
    onError: () => {
      setRunning(false);
    },
    onMutate: (message) => {
      setRunning(true);
      setUserMessage('');
      setMobileView('Chat');

      thread?.messagesList.push({
        id: Date.now().toString(),
        role: 1,
        text: message,
        mentionsList: [],
      });
    },
  });

  const handleOnSubmitUserMessage = useCallback(() => {
    addMessage({ message: userMessage });
  }, [addMessage, userMessage]);

  const handleCopyAsMarkdown = useCallback(() => {
    navigator.clipboard.writeText('');
  }, []);

  const scrollToBottom = useCallback(() => {
    if (!chatContainerRef.current) return;

    observerRef.current?.scrollIntoView({
      behavior: 'smooth',
      block: 'start',
    });
  }, []);

  useEffect(() => {
    if (!thread) return;
    setNowDisplayedStep(thread.currentStepSeqNo);
  }, [thread]);

  useEffect(() => {
    const observerTarget = observerRef.current;
    if (!observerTarget || !chatContainerRef.current) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsAtBottom(entry.isIntersecting);
      },
      {
        root: chatContainerRef.current,
        threshold: 0.1,
      },
    );

    observer.observe(observerTarget);

    return () => {
      observer.disconnect();
    };
  }, []);

  return (
    <div className="flex w-full flex-row gap-[0.875rem] bg-[#F7F7F7] lg:pr-[0.875rem]">
      <div className="shadow-view hidden flex-col items-center gap-2 border border-[#E5E7EB] bg-white py-9 lg:flex px-4 gap-y-7">
        <button onClick={() => router.back()}>
          <ChevronLeft />
        </button>
        <div>
          <span className="text-sm font-bold">Agents</span>
        </div>
        <div className="flex flex-grow overflow-y-auto flex-col p-0.5 gap-y-2 scrollbar-hide">
          <AgentProfileList agentWorks={agentWorks} />
        </div>
      </div>

      <div className="flex flex-grow h-full lg:w-fit flex-col border border-[#E5E7EB] bg-white lg:max-w-[30.25rem]">
        <div className="relative flex items-center justify-center px-14 py-5 lg:py-6">
          <button
            className="absolute left-2 flex lg:hidden"
            onClick={() => router.back()}
          >
            <ChevronLeft className="text-black" />
          </button>

          {mission ? (
            <span className="line-clamp-1 text-sm font-normal">
              {mission.name}
            </span>
          ) : (
            <LoadingSpinner className="m-auto flex h-12 w-12" />
          )}
        </div>
        <div className="flex w-full justify-between px-1.5 gap-x-6 pb-1 lg:hidden overflow-x-auto flex-shrink-0 scrollbar-hide">
          <AgentProfileList agentWorks={agentWorks} />
        </div>

        <div
          ref={chatContainerRef}
          className={classNames(
            'relative h-full flex-col overflow-y-auto border-t border-[#E2E8F0] px-4 py-4',
            {
              'flex flex-grow': mobileView === 'Chat',
              'hidden lg:flex': mobileView !== 'Chat',
            },
          )}
        >
          <ChatSection thread={thread} agentWorks={agentWorks} />
          <div ref={observerRef} className="h-0.5" />

          <div className="sticky bottom-0 flex justify-center">
            <button
              onClick={scrollToBottom}
              className={classNames(
                'flex size-7 items-center justify-center rounded-full bg-gray-400 text-white transition-colors hover:bg-gray-500 lg:size-12',
                {
                  hidden: isAtBottom,
                },
              )}
            >
              <ChevronDown className='size-4 lg:size-7' />
            </button>
          </div>
        </div>

        <div
          className={classNames(
            'h-full flex-col overflow-y-auto border-t border-[#E2E8F0] p-6',
            {
              hidden: mobileView !== 'Workflow',
              'flex flex-grow lg:hidden lg:flex-auto': mobileView === 'Workflow',
            },
          )}
        >
          <WorkflowSection
            loading={isActionWorksLoading}
            stepsList={mission?.stepsList ?? []}
            actionWorks={actionWorks}
            currentStepSeqNo={thread?.currentStepSeqNo ?? 0}
            nowDisplayedStep={nowDisplayedStep ?? 0}
            setNowDisplayedStep={setNowDisplayedStep}
          />
        </div>

        <div
          className={classNames(
            'h-full flex-col overflow-y-auto border-t border-[#E2E8F0] p-6',
            {
              hidden: mobileView !== 'Outcome',
              'flex flex-grow lg:hidden lg:flex-auto': mobileView === 'Outcome',
            },
          )}
        >
          <ResultSection
            markdownContent={''}
            resultCollapsed={resultCollapsed}
            setResultCollapsed={setResultCollapsed}
            onClickCopyButton={handleCopyAsMarkdown}
          />
        </div>

        <MobileNavbar mobileView={mobileView} setMobileView={setMobileView} />

        <div className="relative flex w-full flex-grow items-end px-6 pb-6 lg:pb-[1.875rem]">
          <UserMessageInput
            loading={isRunning || !thread}
            value={userMessage}
            onChange={(v) => setUserMessage(v)}
            onSubmit={handleOnSubmitUserMessage}
          />
        </div>
      </div>

      <div className="hidden h-screen w-full max-w-[calc(100%-8.125rem-30.25rem)] flex-col items-center gap-[0.875rem] pb-[0.875rem] pt-[0.875rem] lg:flex">
        <div className="shadow-card flex max-h-[60%] w-full flex-col overflow-y-auto rounded-[1.25rem] border border-[#E5E7EB] bg-white px-[2.875rem] pb-[2.875rem] pt-8">
          <WorkflowSection
            loading={isActionWorksLoading}
            stepsList={mission?.stepsList ?? []}
            actionWorks={actionWorks}
            currentStepSeqNo={thread?.currentStepSeqNo ?? 0}
            nowDisplayedStep={nowDisplayedStep ?? 0}
            setNowDisplayedStep={setNowDisplayedStep}
          />
        </div>
        <div
          className={classNames(
            'shadow-card flex w-full flex-1 rounded-[1.25rem] border border-[#E5E7EB] bg-white p-4',
            {
              'lg:overflow-hidden': resultCollapsed,
            },
          )}
        >
          <ResultSection
            markdownContent={thread?.result ?? ''}
            resultCollapsed={resultCollapsed}
            setResultCollapsed={setResultCollapsed}
            onClickCopyButton={handleCopyAsMarkdown}
          />
        </div>
      </div>
    </div>
  );
}

function AgentProfileList({
  agentWorks,
}: {
  agentWorks: {
    agent: Agent.AsObject;
    status: AgentWork.Status;
  }[];
}) {
  return (
    <>
      {agentWorks.length > 0 ? (
        agentWorks.map(({ agent, status }) => (
          <AgentProfile
            key={`workflow-agent-${agent.id}`}
            name={agent.name}
            imageUrl={agent.iconUrl}
            imageClassName="p-2 w-14 h-14"
            status={status}
          >
            <AgentProfile.Label className="text-sm font-light">
              {agent.name}
            </AgentProfile.Label>
          </AgentProfile>
        ))
      ) : (
        <LoadingSpinner className="m-auto flex h-12 w-12" />
      )}
    </>
  );
}
