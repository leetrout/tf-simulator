// Static mock data backing the mockup layout. Real data wiring comes in later
// milestones (sim API + websocket); for now this drives the visual shell.

export type TabId = 'sim' | 'scenarios' | 'settings';

export interface Tab {
	id: TabId;
	label: string;
	badge?: string;
}

export const tabs: Tab[] = [
	{ id: 'sim', label: 'Sim' },
	{ id: 'scenarios', label: 'Scenarios', badge: '5' },
	{ id: 'settings', label: 'Settings' }
];

// --- status color presets -------------------------------------------------

const STATUS = {
	sync: 'border-emerald-500/30 bg-emerald-500/10 text-emerald-400',
	untracked: 'border-sky-500/30 bg-sky-500/10 text-sky-400',
	phantom: 'border-rose-500/30 bg-rose-500/10 text-rose-400',
	drift: 'border-amber-500/30 bg-amber-500/10 text-amber-400',
	moved: 'border-violet-500/30 bg-violet-500/10 text-violet-400'
};

// --- graph legend ---------------------------------------------------------

export interface LegendItem {
	label: string;
	dot: string;
}

export const legend: LegendItem[] = [
	{ label: 'in sync', dot: 'bg-emerald-400' },
	{ label: 'untracked', dot: 'bg-sky-400' },
	{ label: 'phantom', dot: 'bg-rose-400' },
	{ label: 'drift', dot: 'bg-amber-400' },
	{ label: 'moved', dot: 'bg-violet-400' }
];

// --- resource graph nodes -------------------------------------------------

export interface GraphNode {
	kind: string;
	name: string;
	status: string;
	statusClass: string;
	wide?: boolean;
}

export const graphNodes: GraphNode[] = [
	{
		kind: 'NETWORK',
		name: 'nimbus_network.main',
		status: 'IN SYNC',
		statusClass: STATUS.sync,
		wide: true
	},
	{ kind: 'SUBNET', name: 'nimbus_subnet.public', status: 'IN SYNC', statusClass: STATUS.sync },
	{ kind: 'SUBNET', name: 'nimbus_subnet.private', status: 'IN SYNC', statusClass: STATUS.sync },
	{
		kind: 'INTERFACE',
		name: 'nimbus_network_interface.web0',
		status: 'IN SYNC',
		statusClass: STATUS.sync
	},
	{
		kind: 'BUCKET',
		name: 'nimbus_storage_bucket.logs',
		status: 'IN SYNC',
		statusClass: STATUS.sync
	},
	{
		kind: 'VM INSTANCE',
		name: 'nimbus_vm_instance.web',
		status: 'IN SYNC',
		statusClass: STATUS.sync
	},
	{
		kind: 'BUCKET',
		name: 'nimbus_storage_bucket.assets',
		status: 'PHANTOM',
		statusClass: STATUS.phantom
	}
];

// --- terraform.tfstate JSON preview ---------------------------------------

export const stateJson = `{
  "version": 4,
  "resources": {
    "nimbus_network.main": {
      "id": "net-1a20",
      "cidr_block": "10.0.0.0/16",
      "region": "us-west-1"
    },
    "nimbus_subnet.public": {
      "id": "subnet-7c1",
      "cidr_block": "10.0.1.0/24",
      "network_id": "net-1a20"
    },
    "nimbus_subnet.private": {
      "id": "subnet-9f2",
      "cidr_block": "10.0.2.0/24",
      "network_id": "net-1a20"
    },
    "nimbus_network_interface.web0": {
      "id": "nic-4d8",
      "subnet_id": "subnet-7c1",
      "private_ip": "10.0.1.20"
    },
    "nimbus_vm_instance.web": {
      "id": "vm-22e",
      "machine_type": "n2-standard-2",
      "nic_id": "nic-4d8"
    }
  }
}`;

// --- scenarios ------------------------------------------------------------

export interface Scenario {
	title: string;
	dot: string;
	badge: string;
	badgeClass: string;
	description: string;
	affects: string;
	command: string;
	loaded?: boolean;
}

export const scenarios: Scenario[] = [
	{
		title: 'In Sync',
		dot: 'bg-emerald-400',
		badge: 'BASELINE',
		badgeClass: STATUS.sync,
		description: 'All 7 resources are tracked and current. terraform plan reports no changes.',
		affects: '— none —',
		command: 'terraform plan = 0 to add, 0 to change, 0 to destroy'
	},
	{
		title: 'Missing in State',
		dot: 'bg-sky-400',
		badge: 'UNTRACKED',
		badgeClass: STATUS.untracked,
		description:
			'nimbus_storage_bucket.logs was created outside Terraform — it is unmanaged until imported.',
		affects: 'nimbus_storage_bucket.logs',
		command: 'terraform import nimbus_storage_bucket.logs logs-bkt'
	},
	{
		title: 'Phantom in State',
		dot: 'bg-rose-400',
		badge: 'DEPLOYED',
		badgeClass: STATUS.phantom,
		description:
			'nimbus_storage_bucket.assets is in state but was deleted in the cloud. A plan will try to recreate it.',
		affects: 'nimbus_storage_bucket.assets',
		command: 'terraform state rm nimbus_storage_bucket.assets',
		loaded: true
	},
	{
		title: 'Config Drift',
		dot: 'bg-amber-400',
		badge: 'DRIFT',
		badgeClass: STATUS.drift,
		description:
			'nimbus_subnet.private.cidr_block changed in the cloud. Applying your config would overwrite it back.',
		affects: 'nimbus_subnet.private',
		command: 'terraform plan -refresh-only # review drift'
	},
	{
		title: 'Renamed / Moved',
		dot: 'bg-violet-400',
		badge: 'MOVED',
		badgeClass: STATUS.moved,
		description:
			'nimbus_network_interface.web0 was renamed to web_primary. Without a moved block it is destroyed and recreated.',
		affects: 'nimbus_network_interface.web0',
		command: 'moved { from = ...web0 to = ...web_primary }'
	}
];

// --- settings -------------------------------------------------------------

export interface SettingRow {
	label: string;
	description: string;
	kind: 'pill' | 'color' | 'toggle';
	value: string | boolean;
	mono?: boolean;
}

export const environmentSettings: SettingRow[] = [
	{
		label: 'Cloud provider',
		description: 'fictional provider for the sim',
		kind: 'pill',
		value: 'Nimbus Cloud'
	},
	{
		label: 'Region',
		description: 'default region for new resources',
		kind: 'pill',
		value: 'us-west-1',
		mono: true
	},
	{
		label: 'State backend',
		description: 'where terraform.tfstate lives',
		kind: 'pill',
		value: 'local',
		mono: true
	},
	{
		label: 'Terraform version',
		description: 'reported in the status bar',
		kind: 'pill',
		value: 'v1.9.2',
		mono: true
	}
];

export const simulatorSettings: SettingRow[] = [
	{
		label: 'Default scenario',
		description: 'loaded on first launch',
		kind: 'pill',
		value: 'In Sync'
	},
	{ label: 'Accent color', description: 'UI highlight color', kind: 'color', value: '#80c5f6' },
	{
		label: 'Dependency edges',
		description: 'dashed links between resources',
		kind: 'toggle',
		value: true
	},
	{
		label: 'Inline drift comments',
		description: 'annotations in the JSON explorer',
		kind: 'toggle',
		value: true
	}
];
