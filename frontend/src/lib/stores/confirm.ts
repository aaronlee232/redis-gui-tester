import { writable } from 'svelte/store';

export interface ConfirmOptions {
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  confirmClass?: string;
  onConfirm: () => void | Promise<void>;
}

interface ConfirmState {
  show: boolean;
  title: string;
  message: string;
  confirmLabel: string;
  cancelLabel: string;
  confirmClass: string;
}

function createConfirmStore() {
  const { subscribe, set, update } = writable<ConfirmState>({
    show: false,
    title: '',
    message: '',
    confirmLabel: 'Confirm',
    cancelLabel: 'Cancel',
    confirmClass: 'btn-primary',
  });

  // Stored outside writable so runConfirm can access it synchronously.
  let onConfirmCallback: (() => void | Promise<void>) | null = null;

  return {
    subscribe,
    open(opts: ConfirmOptions) {
      onConfirmCallback = opts.onConfirm;
      set({
        show: true,
        title: opts.title,
        message: opts.message,
        confirmLabel: opts.confirmLabel ?? 'Confirm',
        cancelLabel: opts.cancelLabel ?? 'Cancel',
        confirmClass: opts.confirmClass ?? 'btn-primary',
      });
    },
    close() {
      update((s) => ({ ...s, show: false }));
    },
    async runConfirm() {
      if (onConfirmCallback) await Promise.resolve(onConfirmCallback());
      update((s) => ({ ...s, show: false }));
    },
  };
}

export const confirm = createConfirmStore();
