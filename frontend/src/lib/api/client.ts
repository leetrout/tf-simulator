// HTTP client for the statesim backend. Paths are same-origin relative so the
// embedded production build "just works" (the Go binary serves both the UI and
// the API from one origin).
import type { GraphData, LoadResult, ScenarioInfo, Settings, Snapshot } from '$lib/sim/types';

const BASE = '';

async function getJSON<T>(path: string): Promise<T> {
	const res = await fetch(BASE + path);
	if (!res.ok) throw new Error(`GET ${path} → ${res.status}`);
	return res.json() as Promise<T>;
}

async function postJSON<T>(path: string): Promise<T> {
	const res = await fetch(BASE + path, { method: 'POST' });
	if (!res.ok) throw new Error(`POST ${path} → ${res.status}`);
	return res.json() as Promise<T>;
}

export const getSnapshot = () => getJSON<Snapshot>('/api/sim/snapshot');
export const getGraph = (view = '') =>
	getJSON<GraphData>('/api/sim/graph' + (view ? `?view=${encodeURIComponent(view)}` : ''));
export const getScenarios = () =>
	getJSON<{ scenarios: ScenarioInfo[] }>('/api/sim/scenarios').then((d) => d.scenarios);
export const getSettings = () => getJSON<Settings>('/api/sim/settings');
export const loadScenario = (id: string) =>
	postJSON<LoadResult>(`/api/sim/scenarios/${encodeURIComponent(id)}/load`);
export const resetSim = () => postJSON<{ ok: boolean }>('/api/sim/reset');

// wsURL derives the websocket endpoint from the current page origin.
export function wsURL(): string {
	if (typeof window === 'undefined') return '';
	const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
	return `${proto}//${window.location.host}/ws`;
}
