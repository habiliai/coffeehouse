import { useAliceApiClient } from '@/hooks/aliceapi';
import { CreateThreadRequest } from '@/proto/aliceapi';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

export function useGetMissions() {
  const client = useAliceApiClient();

  return useQuery({
    queryKey: ['aliceapi.getMissions'] as const,
    initialData: [],
    queryFn: async () => {
      try {
        const { missionsList } = await client
          .getMissions(new Empty())
          .then((res) => res.toObject());

        return missionsList;
      } catch (e) {
        console.error(e);
        throw e;
      }
    },
  });
}

export function useGetAgents() {
  const client = useAliceApiClient();

  return useQuery({
    queryKey: ['aliceapi.getAgents'] as const,
    initialData: [],
    queryFn: async () => {
      try {
        const { agentsList: agents } = await client
          .getAgents(new Empty())
          .then((res) => res.toObject());

        return agents;
      } catch (e) {
        console.error(e);
        throw e;
      }
    },
  });
}

export function useCreateThread({
  onSuccess,
  onError,
}: {
  onSuccess?: (threadId: number) => void;
  onError?: () => void;
} = {}) {
  const client = useAliceApiClient();

  return useMutation({
    mutationKey: ['aliceapi.createThread'] as const,
    async mutationFn({ missionId }: { missionId: number }) {
      const { id: threadId } = await client
        .createThread(new CreateThreadRequest().setMissionId(missionId))
        .then((res) => res.toObject());

      return threadId;
    },
    onSuccess(threadId) {
      onSuccess?.(threadId);
    },
    onError(error) {
      console.error(error);
      onError?.();
    },
  });
}
