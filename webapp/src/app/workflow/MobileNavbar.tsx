import classNames from 'classnames';

export default function MobileNavbar({
  mobileView,
  setMobileView,
}: {
  mobileView: string;
  setMobileView: (view: 'Chat' | 'Workflow' | 'Outcome') => void;
}) {
  return (
    <div className="flex justify-center pb-3.5 lg:hidden gap-x-2">
      <MobileViewButton
        isSelected={mobileView === 'Chat'}
        name="Chat"
        onClick={() => setMobileView('Chat')}
      />
      <MobileViewButton
        isSelected={mobileView === 'Workflow'}
        name="Workflow"
        onClick={() => setMobileView('Workflow')}
      />
      <MobileViewButton
        isSelected={mobileView === 'Outcome'}
        name="Outcome"
        onClick={() => setMobileView('Outcome')}
      />
    </div>
  );
}

function MobileViewButton({
  isSelected,
  name,
  onClick,
}: {
  isSelected: boolean;
  name: string;
  onClick: () => void;
}) {
  return (
    <div
      className={classNames(
        'flex items-center rounded-[6.25rem] ring-2 ring-offset-1',
        {
          'ring-[#4AD15B]': isSelected,
          'ring-transparent': !isSelected,
        },
      )}
    >
      <button
        className={classNames(
          'rounded-[6.25rem] ring-1 px-[1.125rem] py-[0.0625rem] text-sm ring-offset-1',
          {
            'ring-[#4AD15B] bg-[#4AD15B] font-bold text-white': isSelected,
            'ring-[#E5E7EB] font-normal text-black shadow-button':
              !isSelected,
          },
        )}
        onClick={onClick}
      >
        {name}
      </button>
    </div>
  );
}
