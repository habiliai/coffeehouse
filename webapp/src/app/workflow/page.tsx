'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import {
  useAddMessage,
  useGetStatus,
  useGetThread,
  useGetMissionStepStatus,
} from './actions';
import AgentProfile from '@/components/AgentProfile';
import { ChevronLeft } from 'lucide-react';
import { useCallback, useEffect, useMemo, useState } from 'react';
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

  useEffect(() => {
    if (!thread) return;
    setNowDisplayedStep(thread.currentStepSeqNo);
  }, [thread]);

  return (
    <div className="flex w-full flex-row gap-[0.875rem] bg-[#F7F7F7] lg:pr-[0.875rem]">
      <div className="shadow-view hidden w-full max-w-[8.125rem] flex-col items-center gap-2 border border-[#E5E7EB] bg-white py-[2.375rem] lg:flex">
        <button className="mb-14" onClick={() => router.back()}>
          <ChevronLeft />
        </button>
        <div className="mb-7">
          <span className="text-sm font-bold">Agents</span>
        </div>
        <AgentProfileList agentWorks={agentWorks} />
      </div>

      <div className="shadow-view flex h-full w-full flex-col border border-[#E5E7EB] bg-white lg:max-w-[30.25rem]">
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
        <div className="flex w-full justify-center gap-x-6 pb-[0.875rem] lg:hidden">
          <AgentProfileList agentWorks={agentWorks} />
        </div>

        <div
          className={classNames(
            'h-full flex-col gap-4 overflow-y-auto border-t border-[#E2E8F0] px-6 py-9',
            {
              flex: mobileView === 'Chat',
              'hidden lg:flex': mobileView !== 'Chat',
            },
          )}
        >
          <ChatSection thread={thread} agentWorks={agentWorks} />
        </div>

        <div
          className={classNames(
            'h-full flex-col overflow-y-auto border-t border-[#E2E8F0] p-6',
            {
              hidden: mobileView !== 'Workflow',
              'flex lg:hidden': mobileView === 'Workflow',
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
              'flex lg:hidden': mobileView === 'Outcome',
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
            imageClassName="p-2 size-10"
            status={status}
          >
            <AgentProfile.Label className="text-sm font-light">
              {agent.name}
            </AgentProfile.Label>
          </AgentProfile>
        ))
      ) : (
        <LoadingSpinner className="flex h-12 w-12" />
      )}
    </>
  );
}
