<script lang="ts">
  import { fade, fly } from 'svelte/transition';
  import { toast } from '../lib/stores/toast';
</script>

{#if $toast.visible}
  <div
    class="toast toast-end toast-top z-50 p-4"
    in:fly={{ x: 16, duration: 200 }}
    out:fade={{ duration: 150 }}
  >
    <div
      class="alert alert-soft shadow-lg max-w-sm"
      class:alert-success={$toast.type === 'success'}
      class:alert-error={$toast.type === 'error'}
      class:alert-warning={$toast.type === 'warning'}
      class:alert-info={$toast.type === 'info'}
      role="alert"
    >
      <span>{$toast.message}</span>
      <button
        type="button"
        class="btn btn-ghost btn-sm btn-circle"
        aria-label="Dismiss"
        onclick={() => toast.dismiss()}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    </div>
  </div>
{/if}
