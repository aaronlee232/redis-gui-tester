<script lang="ts">
  import { onMount } from 'svelte';
  import { scenarios } from './lib/stores/scenarios';
  import { toast } from './lib/stores/toast';
  import { runAllScenarios } from './lib/api';
  import type { Scenario } from './lib/types';
  import ScenarioList from './components/ScenarioList.svelte';
  import SaveScenarioModal from './components/SaveScenarioModal.svelte';
  import ConfirmModal from './components/ConfirmModal.svelte';
  import Toast from './components/Toast.svelte';

  type Tab = 'all' | 'untested' | 'failed' | 'passed';
  const tabs: Tab[] = ['all', 'untested', 'failed', 'passed'];

  let activeTab = $state<Tab>('all');
  let searchQuery = $state('');
  let showSaveModal = $state(false);
  let editingScenario = $state<Scenario | null>(null);
  let runningAll = $state(false);

  onMount(() => {
    scenarios.fetchAll();
  });

  function openCreateModal() {
    editingScenario = null;
    showSaveModal = true;
  }

  function openEditModal(scenario: Scenario) {
    editingScenario = scenario;
    showSaveModal = true;
  }

  function closeModal() {
    showSaveModal = false;
    editingScenario = null;
  }

  async function handleRunAll() {
    runningAll = true;
    try {
      await runAllScenarios();
      await scenarios.fetchAll();
      toast.show('All scenarios run', 'success');
    } catch {
      toast.show('Run all failed', 'error');
    } finally {
      runningAll = false;
    }
  }
</script>

<div class="min-h-screen bg-base-200 pt-8">
  <main class="container mx-auto max-w-6xl p-4 flex flex-col gap-4">
    <h1 class="text-2xl font-bold">Redis GUI Tester</h1>

    <div class="flex flex-wrap items-center justify-between gap-2">
      <div role="tablist" class="tabs tabs-boxed bg-base-100 p-1 rounded-lg">
        {#each tabs as tab}
          <button
            role="tab"
            class="tab"
            class:tab-active={activeTab === tab}
            aria-selected={activeTab === tab}
            onclick={() => (activeTab = tab)}
          >
            {tab.charAt(0).toUpperCase() + tab.slice(1)}
          </button>
        {/each}
      </div>

      <button type="button" class="btn btn-success" onclick={openCreateModal}>
        New Scenario
      </button>
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <input
        type="search"
        placeholder="Search Scenarios"
        class="input input-bordered flex-1 min-w-0 max-w-md"
        bind:value={searchQuery}
      />
      <button
        type="button"
        class="btn btn-primary"
        onclick={handleRunAll}
        disabled={runningAll}
      >
        {runningAll ? 'Running...' : 'RUN ALL'}
      </button>
    </div>

    <div class="flex-1 min-h-0">
      <ScenarioList {activeTab} {searchQuery} onEdit={openEditModal} />
    </div>
  </main>
</div>

<SaveScenarioModal isOpen={showSaveModal} scenario={editingScenario} onClose={closeModal} />
<ConfirmModal />
<Toast />
