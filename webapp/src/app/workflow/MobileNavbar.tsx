import classNames from 'classnames';

export default function MobileNavbar({
  mobileView,
  setMobileView,
}: {
  mobileView: string;
  setMobileView: (view: 'Chat' | 'Workflow' | 'Outcome') => void;
}) {
  return (
    <div className="flex justify-center pb-[0.875rem] lg:hidden">
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
        'flex items-center rounded-[6.25rem] border-2 p-[1px]',
        {
          'border-[#4AD15B]': isSelected,
          'border-transparent': !isSelected,
        },
      )}
    >
      <button
        className={classNames(
          'rounded-[6.25rem] border px-[1.125rem] py-[0.0625rem] text-sm',
          {
            'border-[#4AD15B] bg-[#4AD15B] font-bold text-white': isSelected,
            'border-[#E5E7EB] font-normal text-black shadow-[0_0_40px_3px_#AEAEAE40]':
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
