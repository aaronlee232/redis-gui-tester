import { writable } from 'svelte/store';
import { getAllScenarios } from '../api';
import type { Scenario } from '../types';

function createScenariosStore() {
  const { subscribe, set } = writable<Scenario[]>([]);

  return {
    subscribe,
    async fetchAll() {
      const data = await getAllScenarios();
      set(data);
    },
  };
}

export const scenarios = createScenariosStore();
