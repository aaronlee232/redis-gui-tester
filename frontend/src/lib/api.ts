import type {
  RunAllScenarioResult,
  RunScenarioResponse,
  Scenario,
  ScenarioPayload,
} from './types';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, init);
  if (!res.ok) {
    const text = await res.text().catch(() => res.statusText);
    throw new Error(`${res.status}: ${text}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export function getAllScenarios(): Promise<Scenario[]> {
  return request<Scenario[]>('/api/scenario/get-all');
}

export function getScenario(id: number): Promise<Scenario> {
  return request<Scenario>(`/api/scenario/get/${id}`);
}

export function createScenario(payload: ScenarioPayload): Promise<void> {
  return request<void>('/api/scenario/create', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
}

export function updateScenario(id: number, payload: ScenarioPayload): Promise<void> {
  return request<void>(`/api/scenario/update/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
}

export function deleteScenario(id: number): Promise<void> {
  return request<void>(`/api/scenario/delete/${id}`, { method: 'DELETE' });
}

export function runScenario(id: number): Promise<RunScenarioResponse> {
  return request<RunScenarioResponse>(`/api/tester/run-scenario/${id}`, { method: 'POST' });
}

export function runAllScenarios(): Promise<RunAllScenarioResult[]> {
  return request<RunAllScenarioResult[]>('/api/tester/run-scenario', { method: 'POST' });
}
