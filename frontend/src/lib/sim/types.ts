// TypeScript mirror of the statesim backend contract (internal/cloud + internal/sim).
// Keep in sync with the Go types when the contract changes.

export type Status = 'in_sync' | 'drift' | 'phantom' | 'untracked' | 'moved';

export interface CloudResource {
	type: string; // nimbus_network, ...
	id: string;
	serial: number;
	attributes: Record<string, string>;
}

export interface StateResource {
	address: string; // nimbus_subnet.public
	module: string;
	type: string;
	name: string;
	id: string;
	attributes: Record<string, string>;
}

export interface AttrDiff {
	field: string;
	state: string;
	cloud: string;
}

export interface ResourceStatus {
	id: string;
	type: string;
	name: string;
	address?: string;
	status: Status;
	detail: string;
	diffs?: AttrDiff[];
}

export interface Counts {
	total: number;
	tracked: number;
	problems: number;
	byStatus: Record<Status, number>;
}

export interface Meta {
	workDir: string;
	stateFile: string;
	stateExists: boolean;
	terraformVersion: string;
	provider: string;
	region: string;
	backend: string;
	loadedScenario?: string;
	lastEvent?: string;
}

export interface Snapshot {
	cloud: CloudResource[];
	state: StateResource[];
	status: ResourceStatus[];
	counts: Counts;
	meta: Meta;
}

export interface GraphNode {
	id: string;
	label: string;
	type: string;
	kind: string;
	name: string;
	status: Status;
	address?: string;
}

export interface GraphEdge {
	source: string;
	target: string;
	label?: string;
	kind?: string; // "phantom" for dashed edges
}

export interface GraphData {
	nodes: GraphNode[];
	edges: GraphEdge[];
}

export interface ScenarioInfo {
	id: string;
	title: string;
	badge: string;
	status: Status;
	description: string;
	affects: string;
	command: string;
	loaded: boolean;
}

export interface LoadResult {
	loaded: string;
	title: string;
	affects: string;
	command: string;
	note?: string;
}

export interface SettingRow {
	label: string;
	description: string;
	kind: 'pill' | 'color' | 'toggle';
	value: string | boolean;
	mono?: boolean;
}

export interface Settings {
	environment: SettingRow[];
	simulator: SettingRow[];
}

export interface WsEvent {
	type: 'cloud.changed' | 'state.changed' | 'scenario.loaded' | 'reset' | 'hello';
	scope: 'cloud' | 'state' | 'sim';
	resourceType?: string;
	id?: string;
	serial?: number;
	detail?: string;
	timestamp: string;
}
