import LoadingSpinner from '@/components/LoadingSpinner';
import UserChatBubble from './UserChatBubble';
import BotChatBubble from './BotChatBubble';
import { Agent, AgentWork, Thread } from '@/proto/habapi';

export default function ChatSection({
  thread,
  agentWorks,
}: {
  thread: Thread.AsObject | null;
  agentWorks: {
    agent: Agent.AsObject;
    status: AgentWork.Status;
  }[];
}) {
  return (
    <div className='flex flex-col gap-4 grow'>
      {!thread ? (
        <LoadingSpinner className="m-auto flex h-12 w-12" />
      ) : (
        <>
          {thread.messagesList.map(
            ({ id, role, text, agent }) => {
              switch (role) {
                case 1:
                  return (
                    <UserChatBubble
                      id={id}
                      key={`thread-user-${id}`}
                      text={text}
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
  );
}
