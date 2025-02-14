import { useHabiliApiClient } from '@/hooks/habapi';
import { CreateThreadRequest } from '@/proto/habapi';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

export function useGetMissions() {
  const client = useHabiliApiClient();

  return useQuery({
    queryKey: ['habapi.getMissions'] as const,
    queryFn: async () => {
      try {
        const missions = await client
          .getMissions(new Empty())
          .then((res) => res.toObject());

        return missions;
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
    queryFn: async () => {
      try {
        const agents = await client
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
  onSuccess?: (threadId: string) => void;
  onError?: () => void;
} = {}) {
  const client = useHabiliApiClient();

  return useMutation({
    mutationKey: ['habapi.createThread'] as const,
    async mutationFn({ missionId }: { missionId: number }) {
      try {
        const result = await client
          .createThread(new CreateThreadRequest().setMissionId(missionId))
          .then((res) => res.toObject());

        return result.id;
      } catch (e) {
        console.error(e);
        throw e;
      }
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
