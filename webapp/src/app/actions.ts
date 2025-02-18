import { useHabiliApiClient } from '@/hooks/habapi';
import { CreateThreadRequest } from '@/proto/habapi';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

export function useGetMissions() {
  const client = useHabiliApiClient();

  return useQuery({
    queryKey: ['habapi.getMissions'] as const,
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
  const client = useHabiliApiClient();

  return useQuery({
    queryKey: ['habapi.getAgents'] as const,
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
  const client = useHabiliApiClient();

  return useMutation({
    mutationKey: ['habapi.createThread'] as const,
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
