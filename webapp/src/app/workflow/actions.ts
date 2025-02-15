import { useMutation, useQuery } from '@tanstack/react-query';

const SAMPLE_AGENTS = [
  { id: 1, name: 'Social Marketer', iconUrl: '' },
  { id: 3, name: 'Developer', iconUrl: '' },
  { id: 5, name: 'Writer', iconUrl: '' },
];
const SAMPLE_THREAD = {
  id: '1',
  messagesList: [
    {
      id: '1',
      role: 1,
      text: 'Hello, I am a user',
      mentionsList: [],
    },
    {
      id: '2',
      role: 2,
      text: 'Hello, I am an assistant',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '3',
      role: 1,
      text: 'How are you?',
      mentionsList: [],
    },
    {
      id: '4',
      role: 2,
      text: 'I am fine, thank you',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '5',
      role: 1,
      text: 'What services do you offer?',
      mentionsList: [],
    },
    {
      id: '6',
      role: 2,
      text: 'I offer digital marketing, research, and development support.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '7',
      role: 1,
      text: 'That sounds great! Can you tell me more about digital marketing?',
      mentionsList: [],
    },
    {
      id: '8',
      role: 2,
      text: 'Sure, digital marketing involves SEO, content creation, social media management, and more.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '9',
      role: 1,
      text: "I'm interested in learning more about SEO.",
      mentionsList: [],
    },
    {
      id: '10',
      role: 2,
      text: 'SEO involves optimizing your website for search engines.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '11',
      role: 1,
      text: 'What about content creation?',
      mentionsList: [],
    },
    {
      id: '12',
      role: 2,
      text: 'Content creation is all about engaging your audience.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '13',
      role: 1,
      text: 'Can you share some case studies?',
      mentionsList: [],
    },
    {
      id: '14',
      role: 2,
      text: 'Sure, I can provide some examples of successful campaigns.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '15',
      role: 1,
      text: "What's the role of social media management?",
      mentionsList: [],
    },
    {
      id: '16',
      role: 2,
      text: 'It involves interactively engaging with customers on various platforms.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '17',
      role: 1,
      text: 'I see. How important is research in marketing?',
      mentionsList: [],
    },
    {
      id: '18',
      role: 2,
      text: 'Research is vital for understanding market trends and audience behavior.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '19',
      role: 1,
      text: 'Do you offer training on these topics?',
      mentionsList: [],
    },
    {
      id: '20',
      role: 2,
      text: 'Yes, we provide training sessions on various digital marketing strategies.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '21',
      role: 1,
      text: 'Could you elaborate on research methods?',
      mentionsList: [],
    },
    {
      id: '22',
      role: 2,
      text: 'Our methods include surveys, interviews, and detailed analytics.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '23',
      role: 1,
      text: "I'm also curious about your digital development support.",
      mentionsList: [],
    },
    {
      id: '24',
      role: 2,
      text: 'That covers a wide range of IT solutions and technical assistance.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '25',
      role: 1,
      text: 'What tools do you use for SEO optimization?',
      mentionsList: [],
    },
    {
      id: '26',
      role: 2,
      text: 'We leverage industry-leading tools and custom analytics platforms.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '27',
      role: 1,
      text: 'How do you measure campaign success?',
      mentionsList: [],
    },
    {
      id: '28',
      role: 2,
      text: 'We focus on metrics like conversions, engagement, and ROI.',
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
    {
      id: '29',
      role: 1,
      text: 'Thanks for all the detailed insights.',
      mentionsList: [],
    },
    {
      id: '30',
      role: 2,
      text: "You're welcome. Feel free to reach out if you have more questions.",
      agent: SAMPLE_AGENTS[0],
      mentionsList: [],
    },
  ],
};

export function useGetThread({ threadId }: { threadId: string }) {
  return useQuery({
    queryKey: ['server.getThread', { threadId }] as const,
    queryFn: () => {
      try {
        // TODO: Fetch missions from the server
        const thread = {
          ...SAMPLE_THREAD,
          title: 'Digital Marketing Support',
          currentStep: 2,
          status: 'done',
          resultContent: '',
          stepsList: [
            {
              id: 1,
              title: 'Step 1',
              status: 'done',
              tasksList: [
                {
                  id: 1,
                  title: 'Step 1 - Task 1',
                  status: 'done',
                  requiredAgents: [SAMPLE_AGENTS[0], SAMPLE_AGENTS[1]],
                },
                {
                  id: 2,
                  title: 'Step 1 - Task 2',
                  status: 'done',
                  requiredAgents: [SAMPLE_AGENTS[2]],
                },
                {
                  id: 3,
                  title: 'Step 1 - Task 3',
                  status: 'done',
                  requiredAgents: [
                    SAMPLE_AGENTS[0],
                    SAMPLE_AGENTS[1],
                    SAMPLE_AGENTS[2],
                  ],
                },
              ],
            },
            {
              id: 2,
              title: 'Step 2',
              status: 'in-progress',
              tasksList: [
                {
                  id: 4,
                  title: 'Step 2 - Task 1',
                  status: 'done',
                  requiredAgents: [SAMPLE_AGENTS[0], SAMPLE_AGENTS[1]],
                },
                {
                  id: 5,
                  title: 'Step 2 - Task 2',
                  status: 'in-progress',
                  requiredAgents: [SAMPLE_AGENTS[2]],
                },
                {
                  id: 6,
                  title: 'Step 2 - Task 3',
                  status: 'pending',
                  requiredAgents: [
                    SAMPLE_AGENTS[0],
                    SAMPLE_AGENTS[1],
                    SAMPLE_AGENTS[2],
                  ],
                },
              ],
            },
            {
              id: 3,
              title: 'Step 3',
              status: 'pending',
              tasksList: [
                {
                  id: 7,
                  title: 'Step 3 - Task 1',
                  status: 'pending',
                  requiredAgents: [SAMPLE_AGENTS[0], SAMPLE_AGENTS[1]],
                },
                {
                  id: 8,
                  title: 'Step 3 - Task 2',
                  status: 'pending',
                  requiredAgents: [SAMPLE_AGENTS[2]],
                },
                {
                  id: 9,
                  title: 'Step 3 - Task 3',
                  status: 'pending',
                  requiredAgents: [
                    SAMPLE_AGENTS[0],
                    SAMPLE_AGENTS[1],
                    SAMPLE_AGENTS[2],
                  ],
                },
              ],
            },
          ],
        };
        return {
          thread,
          agents: SAMPLE_AGENTS,
        };
      } catch (e) {
        console.error(e);
        throw e;
      }
    },
    enabled: !!threadId,
  });
}

export function useAddMessage({
  onSuccess,
  onError,
  onMutate,
}: {
  onSuccess?: () => void;
  onError?: () => void;
  onMutate?: (message: string) => void;
}) {
  return useMutation({
    mutationKey: ['server.addMessage'] as const,
    async mutationFn({ message }: { message: string }) {
      try {
        // TODO: Implement message sending logic
        console.log('User message:', message);
        return null;
      } catch (e) {
        console.error(e);
        throw e;
      }
    },
    onSuccess() {
      onSuccess?.();
    },
    onError(error) {
      console.error(error);
      onError?.();
    },
    onMutate({ message }) {
      onMutate?.(message);
    },
  });
}
