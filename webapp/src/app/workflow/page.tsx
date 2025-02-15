'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useAddMessage, useGetStatus, useGetThread } from './actions';
import AgentProfile from '@/components/AgentProfile';
import { ChevronLeft } from 'lucide-react';
import { useCallback, useMemo, useState } from 'react';
import UserChatBubble from './UserChatBubble';
import BotChatBubble from './BotChatBubble';
import UserMessageInput from './UserMessageInput';
import LoadingSpinner from '@/components/LoadingSpinner';
import { AgentWork } from '@/proto/habapi';

export default function Page() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [isRunning, setRunning] = useState<boolean>(false);
  const [userMessage, setUserMessage] = useState('');

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

      const { text, mentions } = message
        .split(' ')
        .map((word) => {
          const res = { text: word, mentions: [] };
          if (!word.startsWith('@')) {
            return res;
          }

          const mention = word.slice(1);
          return { text: '', mentions: [mention] };
        })
        .reduce((acc, { text, mentions }) => {
          acc.mentions.push(...mentions);
          if (text === '') {
            return acc;
          }

          if (acc.text !== '') {
            acc.text += ' ';
          }
          acc.text += text;
          return acc;
        });

      thread?.messagesList.push({
        id: Date.now().toString(),
        role: 1,
        text: text,
        mentionsList: mentions,
      });
    },
  });

  const handleOnSubmitUserMessage = useCallback(() => {
    addMessage({ message: userMessage });
  }, [addMessage, userMessage]);

  const handleCopyAsMarkdown = useCallback(() => {
    // TODO: Implement copy as markdown logic
  }, []);

  return (
    <div className="flex w-full flex-row gap-[0.875rem] bg-[#F7F7F7] px-6">
      <div className="flex w-full max-w-[8.125rem] flex-col items-center gap-2 border border-[#E5E7EB] bg-white py-[2.375rem] shadow-[0_0_40px_3px_#AEAEAE40]">
        <button className="mb-14" onClick={() => router.back()}>
          <ChevronLeft />
        </button>
        <div className="mb-7">
          <span className="text-sm font-bold">Agents</span>
        </div>
        {agentWorks.length > 0 ? (
          agentWorks.map(({ agent, status }) => (
            <AgentProfile
              key={`workflow-agent-${agent.id}`}
              name={agent.name}
              imageUrl={agent.iconUrl}
              imageClassName="size-10"
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
      </div>

      <div className="flex h-full w-full max-w-[30.25rem] flex-col border border-[#E5E7EB] bg-white shadow-[0_0_40px_3px_#AEAEAE40]">
        <div className="flex justify-center py-6">
          {mission ? (
            <span className="text-sm font-normal">{mission.name}</span>
          ) : (
            <LoadingSpinner className="m-auto flex h-12 w-12" />
          )}
        </div>

        <div className="flex h-full flex-col gap-4 overflow-y-auto border-t border-[#E2E8F0] px-6 py-9">
          {!thread ? (
            <LoadingSpinner className="m-auto flex h-12 w-12" />
          ) : (
            <>
              {thread.messagesList.map(
                ({ id, role, text, agent, mentionsList: mentions }) => {
                  switch (role) {
                    case 1:
                      return (
                        <UserChatBubble
                          key={`thread-user-${id}`}
                          text={text}
                          mentions={mentions}
                        />
                      );
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
                },
              )}
              {agentWorks
                .filter(({ status }) => status === AgentWork.Status.WORKING)
                .map(({ agent }) => (
                  <BotChatBubble
                    key={`thread-bot-completing-${agent.id}`}
                    botName={agent.name}
                    profileImageUrl={agent.iconUrl}
                    working={true}
                  />
                ))}
            </>
          )}
        </div>
        <div className="relative flex w-full flex-grow items-end px-6 pb-[1.875rem]">
          <UserMessageInput
            loading={isRunning || !thread}
            value={userMessage}
            onChange={(v) => setUserMessage(v)}
            onSubmit={handleOnSubmitUserMessage}
          />
        </div>
      </div>

      <div className="flex w-full flex-col gap-[0.875rem] py-[0.875rem]">
        <div className="flex h-[35.5rem] w-full flex-col rounded-[1.25rem] border border-[#E5E7EB] bg-white p-4 shadow-[0_0_40px_3px_#AEAEAE40]">
          {/* TODO: Update the design of the Workflow section */}
          <div className="flex h-full w-full flex-col gap-[0.875rem] px-[1.875rem] py-[1.625rem]">
            <span className="text-[2rem]/[3.375rem] font-bold">Workflow</span>
            {!thread && <LoadingSpinner className="m-auto flex h-12 w-12" />}
            <div className="flex h-full w-full flex-col"></div>
          </div>
        </div>
        <div className="flex h-full w-full flex-col justify-between rounded-[1.25rem] border border-[#E5E7EB] bg-white p-4 shadow-[0_0_40px_3px_#AEAEAE40]">
          <div className="flex h-full w-full flex-col gap-[0.875rem] px-[1.875rem] py-[1.625rem]">
            <span className="text-[2rem]/[3.375rem] font-bold">Result</span>
          </div>
          {/*<Button*/}
          {/*  disabled={thread?.status !== 'done'}*/}
          {/*  onClick={handleCopyAsMarkdown}*/}
          {/*>*/}
          {/*  Copy as Markdown*/}
          {/*</Button>*/}
        </div>
      </div>
    </div>
  );
}
