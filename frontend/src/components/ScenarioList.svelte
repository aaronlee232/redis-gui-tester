<script lang="ts">
  import { scenarios } from '../lib/stores/scenarios';
  import type { Scenario } from '../lib/types';
  import ScenarioItem from './ScenarioItem.svelte';

  let {
    activeTab,
    searchQuery,
    onEdit,
  }: {
    activeTab: string;
    searchQuery: string;
    onEdit: (s: Scenario) => void;
  } = $props();

  let filtered = $derived(
    $scenarios
      .filter((s) => activeTab === 'all' || s.status === activeTab)
      .filter((s) => {
        const q = searchQuery.trim().toLowerCase();
        if (!q) return true;
        return (
          s.title.toLowerCase().includes(q) || s.description.toLowerCase().includes(q)
        );
      }),
  );
</script>

<ul class="space-y-3">
  {#each filtered as scenario (scenario.ID)}
    <li>
      <ScenarioItem {scenario} {onEdit} />
    </li>
  {:else}
    <li class="text-center text-base-content/50 py-12">
      {$scenarios.length === 0
        ? 'No scenarios yet. Create one to get started.'
        : 'No scenarios match the current filter.'}
    </li>
  {/each}
</ul>
