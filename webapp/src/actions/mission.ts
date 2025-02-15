import { useHabiliApiClient } from '@/hooks/habapi';
import { useQuery } from '@tanstack/react-query';
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
