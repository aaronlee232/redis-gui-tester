export type ScenarioStatus = 'untested' | 'passed' | 'failed';

export interface Scenario {
  ID: number;
  title: string;
  description: string;
  commands: string[];
  expected_responses: string[];
  status: ScenarioStatus;
}

export interface ScenarioPayload {
  title: string;
  description: string;
  commands: string[];
}

export interface RunScenarioResponse {
  outputs: string[];
}

export interface RunAllScenarioResult {
  id: number;
  title: string;
  output: string;
  success: boolean;
  error?: string;
}
