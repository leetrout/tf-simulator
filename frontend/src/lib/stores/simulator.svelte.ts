// Reactive simulator store. Holds the live snapshot/graph/scenarios/settings,
// opens the structured websocket, and refetches the affected views on each event.
import * as api from '$lib/api/client';
import type { GraphData, ScenarioInfo, Settings, Snapshot, WsEvent } from '$lib/sim/types';

class SimulatorStore {
	snapshot = $state<Snapshot | null>(null);
	graph = $state<GraphData | null>(null);
	scenarios = $state<ScenarioInfo[]>([]);
	settings = $state<Settings | null>(null);

	connected = $state(false);
	loading = $state(false);
	error = $state<string | null>(null);
	lastEvent = $state<string>('');
	busyScenario = $state<string | null>(null);

	#ws: WebSocket | null = null;
	#reconnectTimer: ReturnType<typeof setTimeout> | null = null;
	#started = false;

	async refresh() {
		this.loading = true;
		try {
			const [snapshot, graph, scenarios, settings] = await Promise.all([
				api.getSnapshot(),
				api.getGraph(),
				api.getScenarios(),
				api.getSettings()
			]);
			this.snapshot = snapshot;
			this.graph = graph;
			this.scenarios = scenarios;
			this.settings = settings;
			this.error = null;
		} catch (e) {
			this.error = e instanceof Error ? e.message : String(e);
		} finally {
			this.loading = false;
		}
	}

	async loadScenario(id: string) {
		this.busyScenario = id;
		try {
			await api.loadScenario(id);
			await this.refresh();
		} catch (e) {
			this.error = e instanceof Error ? e.message : String(e);
		} finally {
			this.busyScenario = null;
		}
	}

	async reset() {
		try {
			await api.resetSim();
			await this.refresh();
		} catch (e) {
			this.error = e instanceof Error ? e.message : String(e);
		}
	}

	// start performs the initial fetch and opens the websocket. Safe to call once
	// from onMount; further calls are no-ops.
	start() {
		if (this.#started) return;
		this.#started = true;
		void this.refresh();
		this.#connect();
	}

	stop() {
		this.#started = false;
		if (this.#reconnectTimer) clearTimeout(this.#reconnectTimer);
		this.#ws?.close();
		this.#ws = null;
	}

	#connect() {
		const url = api.wsURL();
		if (!url) return;
		let ws: WebSocket;
		try {
			ws = new WebSocket(url);
		} catch {
			this.#scheduleReconnect();
			return;
		}
		this.#ws = ws;
		ws.onopen = () => {
			this.connected = true;
		};
		ws.onmessage = (e) => {
			let ev: WsEvent;
			try {
				ev = JSON.parse(e.data);
			} catch {
				return;
			}
			if (ev.type !== 'hello') {
				this.lastEvent = `${ev.type}${ev.detail ? ` · ${ev.detail}` : ''}`;
				void this.refresh();
			}
		};
		ws.onclose = () => {
			this.connected = false;
			this.#scheduleReconnect();
		};
		ws.onerror = () => {
			ws.close();
		};
	}

	#scheduleReconnect() {
		if (!this.#started || this.#reconnectTimer) return;
		this.#reconnectTimer = setTimeout(() => {
			this.#reconnectTimer = null;
			if (this.#started) this.#connect();
		}, 1500);
	}
}

export const sim = new SimulatorStore();
