<script lang="ts">
  import { scenarios } from '../lib/stores/scenarios';
  import { toast } from '../lib/stores/toast';
  import { createScenario, updateScenario } from '../lib/api';
  import type { Scenario } from '../lib/types';

  let {
    isOpen,
    scenario,
    onClose,
  }: {
    isOpen: boolean;
    scenario: Scenario | null;
    onClose: () => void;
  } = $props();

  let title = $state('');
  let description = $state('');
  let commands = $state<string[]>([]);
  let draggedIndex = $state<number | null>(null);
  let submitting = $state(false);

  let dialog: HTMLDialogElement;

  $effect(() => {
    if (!dialog) return;
    if (isOpen) {
      title = scenario?.title ?? '';
      description = scenario?.description ?? '';
      commands = scenario?.commands ? [...scenario.commands] : [];
      if (!dialog.open) dialog.showModal();
    } else {
      if (dialog.open) dialog.close();
    }
  });

  function addCommand() {
    commands = [...commands, ''];
  }

  function removeCommand(index: number) {
    commands = commands.filter((_, i) => i !== index);
  }

  function updateCommand(index: number, value: string) {
    commands = commands.map((cmd, i) => (i === index ? value : cmd));
  }

  function handleDragStart(index: number) {
    draggedIndex = index;
  }

  function handleDragEnd() {
    draggedIndex = null;
  }

  function handleDrop(targetIndex: number) {
    if (draggedIndex === null || draggedIndex === targetIndex) return;
    const next = [...commands];
    const [moved] = next.splice(draggedIndex, 1);
    next.splice(targetIndex, 0, moved);
    commands = next;
    draggedIndex = null;
  }

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    submitting = true;

    const payload = {
      title: title.trim(),
      description: description.trim(),
      commands: commands.filter((c) => c.trim().length > 0),
    };

    try {
      if (scenario?.ID) {
        await updateScenario(scenario.ID, payload);
        toast.show('Scenario saved', 'success');
      } else {
        await createScenario(payload);
        toast.show('Scenario created', 'success');
      }
      await scenarios.fetchAll();
      onClose();
    } catch {
      toast.show('Save failed', 'error');
    } finally {
      submitting = false;
    }
  }
</script>

<dialog
  bind:this={dialog}
  class="modal modal-bottom sm:modal-middle"
  onclose={onClose}
  oncancel={(e) => {
    e.preventDefault();
    onClose();
  }}
>
  <form class="modal-box w-full max-w-2xl" onsubmit={handleSubmit}>
    <h3 class="font-bold text-lg mb-4">
      {scenario ? 'Edit Scenario' : 'Create New Test Scenario'}
    </h3>

    <div class="form-control w-full mb-4">
      <label class="label" for="scenario-title">
        <span class="label-text font-semibold">Scenario Title</span>
      </label>
      <input
        id="scenario-title"
        type="text"
        bind:value={title}
        placeholder="e.g., Basic Key-Value Operations"
        class="input input-bordered w-full"
        required
      />
    </div>

    <div class="form-control w-full mb-4">
      <label class="label" for="scenario-description">
        <span class="label-text font-semibold">Description</span>
      </label>
      <textarea
        id="scenario-description"
        bind:value={description}
        placeholder="Describe what this scenario tests..."
        class="textarea textarea-bordered w-full h-24"
      ></textarea>
    </div>

    <div class="form-control w-full mb-4">
      <div class="label">
        <span class="label-text font-semibold">Commands Sequence</span>
        <span class="label-text-alt text-base-content/50">Drag to reorder.</span>
      </div>

      <div
        role="list"
        class="space-y-2 p-4 bg-base-200 rounded-lg border border-base-300 min-h-24"
        ondragover={(e) => e.preventDefault()}
      >
        {#if commands.length === 0}
          <div class="text-center text-base-content/40 py-6">
            <p>No commands added yet. Click "Add Command" to get started.</p>
          </div>
        {/if}

        {#each commands as cmd, index (index)}
          <div
            role="listitem"
            class="flex gap-2 items-center p-3 bg-base-100 rounded-lg border border-base-300"
            class:opacity-50={draggedIndex === index}
            draggable="true"
            ondragstart={() => handleDragStart(index)}
            ondragend={handleDragEnd}
            ondragover={(e) => e.preventDefault()}
            ondrop={() => handleDrop(index)}
          >
            <div
              class="cursor-move flex items-center justify-center w-8 h-8 text-base-content/40 hover:text-base-content/70 shrink-0"
            >
              <svg fill="currentColor" viewBox="0 0 32 32" class="w-5 h-5">
                <rect x="10" y="6" width="4" height="4" />
                <rect x="18" y="6" width="4" height="4" />
                <rect x="10" y="14" width="4" height="4" />
                <rect x="18" y="14" width="4" height="4" />
                <rect x="10" y="22" width="4" height="4" />
                <rect x="18" y="22" width="4" height="4" />
              </svg>
            </div>

            <span class="text-sm font-mono text-base-content/60 shrink-0">redis-cli</span>

            <input
              type="text"
              value={cmd}
              oninput={(e) => updateCommand(index, (e.target as HTMLInputElement).value)}
              placeholder='SET mykey "Hello world"'
              class="input input-bordered flex-1"
            />

            <button
              type="button"
              class="btn btn-sm btn-ghost text-error shrink-0"
              aria-label="Remove command"
              onclick={() => removeCommand(index)}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-4 w-4"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fill-rule="evenodd"
                  d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>
          </div>
        {/each}
      </div>

      <button
        type="button"
        class="btn btn-sm btn-outline btn-primary mt-3"
        onclick={addCommand}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fill-rule="evenodd"
            d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
            clip-rule="evenodd"
          />
        </svg>
        Add Command
      </button>
    </div>

    <div class="modal-action">
      <button type="button" class="btn btn-ghost" onclick={onClose}>Cancel</button>
      <button type="submit" class="btn btn-primary" disabled={submitting}>
        {submitting ? 'Saving...' : scenario ? 'Save Changes' : 'Create'}
      </button>
    </div>
  </form>

  <form method="dialog" class="modal-backdrop">
    <button type="button" onclick={onClose}>close</button>
  </form>
</dialog>
