<script lang="ts">
  import { onMount } from 'svelte';
  import type { Scenario } from '../lib/types';
  import { scenarios } from '../lib/stores/scenarios';
  import { toast } from '../lib/stores/toast';
  import { confirm } from '../lib/stores/confirm';
  import { deleteScenario as apiDelete, runScenario as apiRun } from '../lib/api';

  let {
    scenario,
    onEdit,
  }: {
    scenario: Scenario;
    onEdit: (s: Scenario) => void;
  } = $props();

  let isOpen = $state(false);
  let runOutputs = $state<string[]>([]);
  let running = $state(false);

  onMount(() => {
    try {
      const raw = localStorage.getItem('openScenarioIds');
      if (raw) {
        const ids: number[] = JSON.parse(raw);
        isOpen = ids.includes(scenario.ID);
      }
    } catch {
      // ignore localStorage errors
    }
  });

  function persistOpenState(id: number, open: boolean) {
    try {
      const raw = localStorage.getItem('openScenarioIds');
      const ids: number[] = raw ? JSON.parse(raw) : [];
      if (open) {
        if (!ids.includes(id)) ids.push(id);
      } else {
        const idx = ids.indexOf(id);
        if (idx !== -1) ids.splice(idx, 1);
      }
      localStorage.setItem('openScenarioIds', JSON.stringify(ids));
    } catch {
      // ignore localStorage errors
    }
  }

  function onToggle(e: Event) {
    isOpen = (e.target as HTMLInputElement).checked;
    persistOpenState(scenario.ID, isOpen);
  }

  function statusIcon(status: Scenario['status']): string {
    switch (status) {
      case 'passed':
        return '✓';
      case 'failed':
        return '✗';
      default:
        return '−';
    }
  }

  function statusColor(status: Scenario['status']): string {
    switch (status) {
      case 'passed':
        return 'text-success';
      case 'failed':
        return 'text-error';
      default:
        return 'text-base-content/40';
    }
  }

  async function handleRun() {
    running = true;
    try {
      const result = await apiRun(scenario.ID);
      runOutputs = result.outputs;
      await scenarios.fetchAll();
      toast.show('Scenario run completed', 'success');
    } catch {
      toast.show('Run failed', 'error');
    } finally {
      running = false;
    }
  }

  function handleDelete() {
    confirm.open({
      title: 'Delete scenario?',
      message: 'This cannot be undone.',
      confirmLabel: 'Delete',
      confirmClass: 'btn-error',
      onConfirm: async () => {
        try {
          await apiDelete(scenario.ID);
          await scenarios.fetchAll();
          toast.show('Scenario deleted', 'success');
        } catch {
          toast.show('Delete failed', 'error');
        }
      },
    });
  }
</script>

<div class="collapse collapse-plus border border-base-300 bg-base-100 rounded-lg">
  <input type="checkbox" class="peer" checked={isOpen} onchange={onToggle} />
  <div
    class="collapse-title flex items-center gap-4 font-semibold text-base peer-checked:bg-base-200 transition-colors peer-checked:font-bold"
  >
    <span class="text-xl shrink-0 {statusColor(scenario.status)}">
      {statusIcon(scenario.status)}
    </span>
    <span class="grow">{scenario.title}</span>
  </div>

  <div class="collapse-content pt-4">
    <div class="divider my-0"></div>

    <div class="flex w-full">
      <!-- Left column: description, commands, actions -->
      <div class="flex-1 flex flex-col space-y-4">
        <div>
          <div class="label pb-2">
            <span class="label-text font-semibold">Description</span>
          </div>
          <p class="text-sm text-base-content/80 leading-relaxed">{scenario.description}</p>
        </div>

        <div>
          <div class="label pb-2">
            <span class="label-text font-semibold">Commands</span>
          </div>
          <div
            class="bg-neutral text-neutral-content p-3 rounded-lg text-xs overflow-x-auto border border-base-300 font-mono"
          >
            <code>
              {#each scenario.commands as cmd}
                <div class="py-1">
                  <span class="opacity-50">$ redis-cli</span>
                  <span> {cmd}</span>
                </div>
              {/each}
            </code>
          </div>
        </div>

        <div class="mt-auto flex flex-row gap-2">
          <div class="grow pt-2">
            <button
              type="button"
              class="min-w-full btn btn-success btn-sm"
              onclick={handleRun}
              disabled={running}
            >
              {running ? 'Running...' : 'Run Scenario'}
            </button>
          </div>
          <div class="pt-2">
            <button
              type="button"
              class="btn btn-primary btn-sm"
              onclick={() => onEdit(scenario)}
            >
              Edit
            </button>
          </div>
          <div class="pt-2">
            <button type="button" class="btn btn-error btn-sm" onclick={handleDelete}>
              Delete
            </button>
          </div>
        </div>
      </div>

      <!-- Vertical divider -->
      <div class="grow-0 divider divider-horizontal"></div>

      <!-- Right column: run outputs + expected responses -->
      <div class="flex-1 space-y-4 flex flex-col">
        <div class="flex-1 flex flex-col">
          <div class="label pb-2">
            <span class="label-text font-semibold">Redis Responses</span>
          </div>
          <div
            class="bg-neutral text-success p-3 rounded-lg font-mono text-xs overflow-y-auto border border-base-300 flex-1 min-h-32"
          >
            {#each runOutputs as output}
              <div class="py-1">{output}</div>
            {:else}
              <span class="opacity-40">Run scenario to see responses.</span>
            {/each}
          </div>
        </div>

        <div class="flex-1 flex flex-col">
          <div class="label pb-2">
            <span class="label-text font-semibold">Expected Responses</span>
          </div>
          <div
            class="bg-neutral text-info p-3 rounded-lg font-mono text-xs overflow-y-auto border border-base-300 flex-1 min-h-32"
          >
            <code>
              {#each scenario.expected_responses as res}
                <div class="py-1">{res}</div>
              {:else}
                <span class="opacity-40">No expected responses recorded yet.</span>
              {/each}
            </code>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
