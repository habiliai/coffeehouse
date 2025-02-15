import { useHabiliApiClient } from '@/hooks/habapi';
import { useQuery } from '@tanstack/react-query';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

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
