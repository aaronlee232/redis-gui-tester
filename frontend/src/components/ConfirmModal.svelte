<script lang="ts">
  import { confirm } from '../lib/stores/confirm';

  let dialog: HTMLDialogElement;

  $effect(() => {
    if (!dialog) return;
    if ($confirm.show) {
      if (!dialog.open) dialog.showModal();
    } else {
      if (dialog.open) dialog.close();
    }
  });
</script>

<dialog
  bind:this={dialog}
  class="modal modal-bottom sm:modal-middle"
  onclose={() => confirm.close()}
  oncancel={(e) => {
    e.preventDefault();
    confirm.close();
  }}
>
  <div class="modal-box">
    <h3 class="font-bold text-lg">{$confirm.title}</h3>
    <p class="py-2 text-base-content/80">{$confirm.message}</p>
    <div class="modal-action">
      <button type="button" class="btn btn-ghost" onclick={() => confirm.close()}>
        {$confirm.cancelLabel}
      </button>
      <button
        type="button"
        class="btn {$confirm.confirmClass}"
        onclick={() => confirm.runConfirm()}
      >
        {$confirm.confirmLabel}
      </button>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button type="button" onclick={() => confirm.close()}>close</button>
  </form>
</dialog>
