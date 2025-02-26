import classNames from 'classnames';

export default function MobileNavbar({
  mobileView,
  setMobileView,
}: {
  mobileView: string;
  setMobileView: (view: 'Chat' | 'Workflow' | 'Outcome') => void;
}) {
  return (
    <div className="flex justify-center gap-x-2 pb-3.5 lg:hidden">
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
          'rounded-[6.25rem] px-[1.125rem] py-[0.0625rem] text-sm ring-1 ring-offset-1',
          {
            'bg-[#4AD15B] font-bold text-white ring-[#4AD15B]': isSelected,
            'font-normal text-black shadow-button ring-[#E5E7EB]': !isSelected,
          },
        )}
        onClick={onClick}
      >
        {name}
      </button>
    </div>
  );
}
