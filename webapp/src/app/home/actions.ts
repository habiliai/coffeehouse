import { useQuery } from '@tanstack/react-query';

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
  agentPreset: number[];
};

const SAMPLE_MISSIONS: Mission[] = [
  { id: 1, name: 'Mission 1', agentPreset: [1, 3, 5] },
  { id: 2, name: 'Mission 2', agentPreset: [2, 4] },
  { id: 3, name: 'Mission 3', agentPreset: [2, 1, 4, 5, 6] },
];

export function useGetMissions() {
  return useQuery({
    queryKey: ['server.getMissions'] as const,
    queryFn: () => {
      try {
        // TODO: Fetch missions from the server
        const missions: Mission[] = SAMPLE_MISSIONS;
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
        const agents: Agent[] = SAMPLE_AGENTS;
        return agents;
      } catch (e) {
        console.error(e);
        throw e;
      }
    },
  });
}
