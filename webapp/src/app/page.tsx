'use client';

import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { motion, LayoutGroup } from 'framer-motion';
import { useCreateThread, useGetAgents, useGetMissions } from './actions';
import AgentProfile from '@/components/AgentProfile';
import { VerticalCarousel } from '@/components/VerticalCarousel';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Agent, Mission } from '@/proto/habapi';
import LoadingSpinner from '@/components/LoadingSpinner';

const DEFAULT_AGENT_SLOT = Array(3).fill(null);
export default function Home() {
  const router = useRouter();
  const { data: missions, isLoading: isMissionsLoading } = useGetMissions();
  const { data: agents, isLoading: isAgentsLoading } = useGetAgents();

  const { mutate: createThread } = useCreateThread({
    onSuccess: (threadId) => {
      router.push(`/workflow?thread_id=${threadId}`);
    },
  });

  const [selectedMission, setSelectedMission] = useState<Mission.AsObject>();
  const [agentSlots, setAgentSlots] =
    useState<(Agent.AsObject | null)[]>(DEFAULT_AGENT_SLOT);

  const handleMissionChange = useCallback(
    (missionId: number) => {
      const mission = missions.find((m) => m.id === missionId);
      if (!mission) return;

      setSelectedMission(mission);
      setAgentSlots(Array(mission.agentIdsList.length).fill(null));
    },
    [missions],
  );

  const handleAgentSlotChange = useCallback(() => {
    if (!selectedMission || !agents) return;

    const matchedAgents = selectedMission.agentIdsList.map(
      (agentId) => agents.find((agent) => agent.id === agentId) ?? null,
    );
    setAgentSlots(matchedAgents);
  }, [agents, selectedMission]);

  const handleCreateThread = useCallback(() => {
    if (!selectedMission) return;
    createThread({ missionId: selectedMission.id });
  }, [selectedMission, createThread]);

  const unassignedAgents = useMemo(() => {
    const assignedAgentIds = agentSlots
      .filter((slot) => slot !== null)
      .map((slot) => slot.id);

    return agents.filter((agent) => !assignedAgentIds.includes(agent.id));
  }, [agents, agentSlots]);

  const isLoading = useMemo(() => {
    return isMissionsLoading || isAgentsLoading;
  }, [isAgentsLoading, isMissionsLoading]);

  useEffect(() => {
    // To ensure smooth transitions for the AgentProfile animations,
    // we trigger handleAgentSlotChange only after agentSlots has been fully cleared.
    // This allows the framer-motion layoutId based animations to work correctly.
    if (!agentSlots.every((slot) => slot === null)) return;
    handleAgentSlotChange();
  }, [agentSlots, handleAgentSlotChange]);

  return (
    <LayoutGroup>
      {isLoading && <LoadingSpinner className="m-auto flex h-24 w-24" />}
      {!isLoading && (
        <div className="flex h-full w-full flex-col items-center gap-y-8 p-6 lg:justify-center">
          <h1 className="text-2xl font-bold md:text-3xl">
            Agent Collaborative Network
          </h1>

          <div className="flex w-full flex-col items-center justify-center gap-x-10 lg:flex-row">
            <h2 className="text-2xl font-bold">Missions</h2>
            <VerticalCarousel
              itemClassName="text-xl lg:text-3xl"
              items={
                missions.map((mission) => ({
                  value: mission.id.toString(),
                  name: mission.name,
                })) ?? []
              }
              selectedValue={selectedMission?.id.toString()}
              onClick={(value) => handleMissionChange(Number(value))}
            />
          </div>

          <div className="flex flex-col items-center gap-4 lg:flex-row lg:items-end">
            <div className="flex flex-wrap justify-center gap-4">
              {agentSlots.map((agent, index) => (
                <div
                  key={`agent-slot-${index}`}
                  className="flex h-32 w-32 items-center justify-center border border-gray-400"
                >
                  {agent && (
                    <motion.div layoutId={`agent-${agent.id}`}>
                      <AgentProfile
                        className="w-32"
                        imageClassName="w-14 h-14"
                        name={agent.name}
                        imageUrl={agent.iconUrl}
                      >
                        <AgentProfile.Label>{agent.name}</AgentProfile.Label>
                      </AgentProfile>
                    </motion.div>
                  )}
                </div>
              ))}
            </div>
            <Button
              className="w-full max-w-80 lg:w-fit"
              onClick={handleCreateThread}
            >
              Mission Start
            </Button>
          </div>
          <div className="flex h-24 flex-wrap justify-center gap-4">
            {unassignedAgents?.map((agent) => (
              <motion.div key={agent.id} layoutId={`agent-${agent.id}`}>
                <AgentProfile
                  className="w-32"
                  imageClassName="w-14 h-14"
                  name={agent.name}
                  imageUrl={agent.iconUrl}
                >
                  <AgentProfile.Label>{agent.name}</AgentProfile.Label>
                </AgentProfile>
              </motion.div>
            ))}
          </div>
        </div>
      )}
    </LayoutGroup>
  );
}
