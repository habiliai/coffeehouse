import { Share } from 'lucide-react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';

export const useNotifyPWA = () => {
  const [deferredPrompt, setDeferredPrompt] =
    useState<BeforeInstallPromptEvent | null>(null);

  const isIos = useMemo((): boolean => {
    return /iphone|ipad|ipod/.test(window.navigator.userAgent.toLowerCase());
  }, []);

  const isSafari = useMemo((): boolean => {
    const ua = window.navigator.userAgent;
    return ua.includes('Safari') && !ua.includes('Chrome');
  }, []);

  const isPWAQuery = useMemo((): boolean => {
    const params = new URLSearchParams(window.location.search);
    return params.get('is_pwa') === 'true';
  }, []);

  useEffect(() => {
    const handleBeforeInstallPrompt = (e: Event) => {
      e.preventDefault();
      setDeferredPrompt(e as BeforeInstallPromptEvent);
    };

    window.addEventListener(
      'beforeinstallprompt',
      handleBeforeInstallPrompt as EventListener,
    );
    return () => {
      window.removeEventListener(
        'beforeinstallprompt',
        handleBeforeInstallPrompt as EventListener,
      );
    };
  }, []);

  const promptInstall = useCallback((): void => {
    if (isPWAQuery) return;

    if (deferredPrompt) {
      toast('Install app', {
        description: (
          <span className="text-black">
            To get back here quickly, install Alice
          </span>
        ),
        action: {
          label: 'Install',
          onClick: async () => {
            await deferredPrompt.prompt();
            setDeferredPrompt(null);
          },
        },
      });
    } else if (
      isIos ||
      (isSafari && window.navigator.platform.includes('Mac'))
    ) {
      toast('Install app', {
        description: (
          <div className="flex flex-col text-black">
            <span>To get back here quickly,</span>
            <span className="flex items-center">
              tap <Share height={18} width={18} /> then &apos;Add to Home
              Screen&apos;
            </span>
          </div>
        ),
        action: {
          label: 'OK',
          onClick: () => {},
        },
      });
    }

    // The PWA installation feature is not available in this environment.
    return;
  }, [deferredPrompt, isIos, isSafari, isPWAQuery]);

  return { promptInstall };
};
