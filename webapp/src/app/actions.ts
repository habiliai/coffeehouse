import { useMutation, useQuery } from '@tanstack/react-query';

export type Agent = {
  id: number;
  name: string;
  iconUrl: string;
};

const SAMPLE_AGENTS: Agent[] = [
  { id: 1, name: 'Social Marketer', iconUrl: '' },
  { id: 2, name: 'Researcher', iconUrl: '' },
  { id: 3, name: 'Developer', iconUrl: '' },
  { id: 4, name: 'Designer', iconUrl: '' },
  { id: 5, name: 'Writer', iconUrl: '' },
  { id: 6, name: 'Analyst', iconUrl: '' },
];

export type Mission = {
  id: number;
  name: string;
  agentsList: Agent[];
};

const SAMPLE_MISSIONS: Mission[] = [
  {
    id: 1,
    name: 'Mission 1',
    agentsList: [SAMPLE_AGENTS[0], SAMPLE_AGENTS[2], SAMPLE_AGENTS[4]],
  },
  {
    id: 2,
    name: 'Mission 2',
    agentsList: [SAMPLE_AGENTS[1], SAMPLE_AGENTS[3]],
  },
  {
    id: 3,
    name: 'Mission 3',
    agentsList: [
      SAMPLE_AGENTS[1],
      SAMPLE_AGENTS[0],
      SAMPLE_AGENTS[3],
      SAMPLE_AGENTS[4],
      SAMPLE_AGENTS[5],
    ],
  },
  {
    id: 4,
    name: 'Mission 4',
    agentsList: [SAMPLE_AGENTS[2], SAMPLE_AGENTS[4], SAMPLE_AGENTS[5]],
  },
];

export function useGetMissions() {
  return useQuery({
    queryKey: ['server.getMissions'] as const,
    queryFn: () => {
      try {
        // TODO: Fetch missions from the server
        const missions = SAMPLE_MISSIONS;
        return missions;
      } catch (e) {
        console.error(e);
        throw e;
      }
    },
  });
}

export function useGetAgents() {
  return useQuery({
    queryKey: ['server.getAgents'] as const,
    queryFn: () => {
      try {
        // TODO: Fetch agents from the server
        const agents = SAMPLE_AGENTS;
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
  return useMutation({
    mutationKey: [''] as const,
    async mutationFn() {
      try {
        // TODO: Implement thread creation logic
        return '1234';
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
