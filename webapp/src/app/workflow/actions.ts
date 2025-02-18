import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useHabiliApiClient } from '@/hooks/habapi';
import {
  AddMessageRequest,
  GetThreadStatusRequest,
  MissionId,
  Thread,
  ThreadId,
} from '@/proto/habapi';
import { useState } from 'react';

export function useGetThread({ threadId }: { threadId: number | null }) {
  const habapi = useHabiliApiClient();
  return useQuery({
    refetchOnWindowFocus: true,
    enabled: !!threadId,
    queryKey: ['server.getThread', { threadId, habapi }] as const,
    initialData: {
      thread: null,
      mission: null,
      lastMessageId: '',
    },
    queryFn: async ({ queryKey: [{}, { threadId, habapi }] }) => {
      try {
        if (!threadId) {
          return { thread: null, mission: null, lastMessageId: '' };
        }

        const thread = await habapi
          .getThread(new ThreadId().setId(threadId))
          .then((r) => r.toObject())
          .then((r) => {
            return { currentStep: 3, status: 'done', ...r };
          });

        const mission = await habapi
          .getMission(new MissionId().setId(thread.missionId))
          .then((r) => r.toObject());

        const lastMessageId = thread.messagesList.slice(-1)[0]?.id ?? '';

        return {
          thread,
          mission,
          lastMessageId,
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
  const habapi = useHabiliApiClient();
  return useMutation({
    mutationKey: ['server.addMessage', { threadId }] as const,
    mutationFn: async ({ message }: { message: string }) => {
      if (!threadId) {
        throw new Error('Thread not found');
      }
      await habapi.addMessage(
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
  lastMessageId: string;
}) {
  const habapi = useHabiliApiClient();
  const queryClient = useQueryClient();
  return useQuery({
    refetchInterval: 1000,
    refetchIntervalInBackground: true,
    enabled: !!threadId,
    queryKey: ['server.getAgentsStatus', { habapi }] as const,
    initialData: {
      agentWorks: [],
      hasNewMessage: false,
    },
    queryFn: async ({ queryKey: [{}, { habapi }] }) => {
      try {
        if (!threadId) {
          throw new Error('Thread not found');
        }

        const [agentWorks, hasNewMessage] = await Promise.all([
          habapi.getAgentsStatus(new ThreadId().setId(threadId)).then((r) => {
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
          habapi
            .getThreadStatus(
              new GetThreadStatusRequest()
                .setThreadId(threadId)
                .setLastMessageId(lastMessageId),
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
