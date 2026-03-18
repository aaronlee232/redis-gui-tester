import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

interface ToastState {
  visible: boolean;
  message: string;
  type: ToastType;
}

function createToastStore() {
  const { subscribe, set, update } = writable<ToastState>({
    visible: false,
    message: '',
    type: 'info',
  });

  let timeoutId: ReturnType<typeof setTimeout> | undefined;

  return {
    subscribe,
    show(message: string, type: ToastType = 'info', duration = 4000) {
      if (timeoutId) clearTimeout(timeoutId);
      set({ visible: true, message, type });
      timeoutId = setTimeout(() => this.dismiss(), duration);
    },
    dismiss() {
      if (timeoutId) {
        clearTimeout(timeoutId);
        timeoutId = undefined;
      }
      update((s) => ({ ...s, visible: false }));
    },
  };
}

export const toast = createToastStore();
