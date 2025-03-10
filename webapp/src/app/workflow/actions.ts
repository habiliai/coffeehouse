import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useAliceApiClient } from '@/hooks/aliceapi';
import {
  AddMessageRequest,
  GetMissionStepStatusRequest,
  GetThreadStatusRequest,
  MissionId,
  ThreadId,
} from '@/proto/aliceapi';

export function useGetThread({ threadId }: { threadId: number | null }) {
  const aliceapi = useAliceApiClient();
  return useQuery({
    refetchOnWindowFocus: true,
    enabled: !!threadId,
    queryKey: ['server.getThread', { threadId, aliceapi }] as const,
    initialData: {
      thread: null,
      mission: null,
      lastMessageId: 0,
      actionWorks: [],
    },
    queryFn: async ({ queryKey: [{}, { threadId, aliceapi }] }) => {
      try {
        if (!threadId) {
          return { thread: null, mission: null, lastMessageId: 0, actionWorks: [] };
        }

        const thread = await aliceapi
          .getThread(new ThreadId().setId(threadId))
          .then((r) => r.toObject())
          .then((r) => {
            return { status: 'done', ...r };
          });

        const mission = await aliceapi
          .getMission(new MissionId().setId(thread.missionId))
          .then((r) => r.toObject());

        const lastMessageId = thread.messagesList.slice(-1)[0]?.id ?? '';

        return {
          thread,
          mission,
          lastMessageId,
          actionWorks: thread.actionWorksList,
        };
      } catch (e) {
        console.error(e);
        throw e;
      }
    },
  });
}

export function useAddMessage({
  threadId,
  onSuccess,
  onError,
  onMutate,
}: {
  threadId: number | null;
  onSuccess?: () => void;
  onError?: () => void;
  onMutate?: (message: string) => void;
}) {
  const queryClient = useQueryClient();
  const aliceapi = useAliceApiClient();
  return useMutation({
    mutationKey: ['server.addMessage', { threadId }] as const,
    mutationFn: async ({ message }: { message: string }) => {
      if (!threadId) {
        throw new Error('Thread not found');
      }
      await aliceapi.addMessage(
        new AddMessageRequest().setThreadId(threadId).setMessage(message),
      );
      return null;
    },
    onSuccess: async () => {
      onSuccess?.();
    },
    onError: async (error) => {
      console.error(error);
      onError?.();
    },
    onSettled: async () => {
      await queryClient.invalidateQueries({ queryKey: ['server.getThread'] });
    },
    onMutate: async ({ message }) => {
      onMutate?.(message);
    },
  });
}

export function useGetStatus({
  threadId,
  lastMessageId,
}: {
  threadId: number | null;
  lastMessageId: number;
}) {
  const aliceapi = useAliceApiClient();
  const queryClient = useQueryClient();
  return useQuery({
    refetchInterval: 1000,
    refetchIntervalInBackground: true,
    enabled: !!threadId,
    queryKey: ['server.getAgentsStatus', { aliceapi }] as const,
    initialData: {
      agentWorks: [],
      hasNewMessage: false,
    },
    queryFn: async ({ queryKey: [{}, { aliceapi }] }) => {
      try {
        if (!threadId) {
          throw new Error('Thread not found');
        }

        const [agentWorks, hasNewMessage] = await Promise.all([
          aliceapi.getAgentsStatus(new ThreadId().setId(threadId)).then((r) => {
            const { worksList: works } = r.toObject();
            const res = works.map((w) => {
              if (!w.agent) {
                throw new Error('Agent not found');
              }
              return {
                agent: w.agent,
                status: w.status,
              };
            });
            res.sort((a, b) => a.agent.id - b.agent.id);

            return res;
          }),
          aliceapi
            .getThreadStatus(
              new GetThreadStatusRequest()
                .setThreadId(threadId)
                .setLastMessageId(lastMessageId)
            )
            .then(async (r) => {
              const { hasNewMessage } = r.toObject();
              if (hasNewMessage) {
                await queryClient.invalidateQueries({
                  queryKey: ['server.getThread'],
                });
              }
              return hasNewMessage;
            }),
        ]);

        return {
          agentWorks,
          hasNewMessage,
        };
      } catch (e) {
        console.error(e);
        throw e;
      }
    },
  });
}

export function useGetMissionStepStatus({
  threadId,
  selectedSeqNo,
}: {
  threadId: number | null;
  selectedSeqNo: number | null;
}) {
  const aliceapi = useAliceApiClient();
  return useQuery({
    refetchInterval: 1000,
    refetchIntervalInBackground: true,
    enabled: !!threadId && !!selectedSeqNo,
    queryKey: [
      'server.getMissionStepStatus',
      { threadId, selectedSeqNo, aliceapi },
    ] as const,
    initialData: {
      initialLoading: true,
      actionWorks: [],
    },
    queryFn: async ({
      queryKey: [{}, { threadId, selectedSeqNo, aliceapi }],
    }) => {
      try {
        if (!threadId || !selectedSeqNo) {
          return {
            initialLoading: false,
            actionWorks: [],
          };
        }

        const works = await aliceapi
          .getMissionStepStatus(
            new GetMissionStepStatusRequest()
              .setThreadId(threadId)
              .setStepSeqNo(selectedSeqNo),
          )
          .then((r) => r.toObject());

        return {
          initialLoading: false,
          actionWorks: works.actionWorksList ?? [],
        };
      } catch (e) {
        console.error(e);
        throw e;
      }
    },
  });
}
