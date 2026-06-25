// Shared status presentation: maps the backend Status vocabulary to the UI
// legend colors (emerald/sky/rose/amber/violet) used across panes.
import type { Status } from './types';

export interface StatusMeta {
	label: string;
	badge: string; // border/bg/text classes for a badge
	dot: string; // background class for a legend dot
	text: string; // text color class
	hex: string; // canvas color (matches the Tailwind *-400 dot)
}

export const STATUS_META: Record<Status, StatusMeta> = {
	in_sync: {
		label: 'in sync',
		badge: 'border-emerald-500/30 bg-emerald-500/10 text-emerald-400',
		dot: 'bg-emerald-400',
		text: 'text-emerald-400',
		hex: '#34d399'
	},
	untracked: {
		label: 'untracked',
		badge: 'border-sky-500/30 bg-sky-500/10 text-sky-400',
		dot: 'bg-sky-400',
		text: 'text-sky-400',
		hex: '#38bdf8'
	},
	phantom: {
		label: 'phantom',
		badge: 'border-rose-500/30 bg-rose-500/10 text-rose-400',
		dot: 'bg-rose-400',
		text: 'text-rose-400',
		hex: '#fb7185'
	},
	drift: {
		label: 'drift',
		badge: 'border-amber-500/30 bg-amber-500/10 text-amber-400',
		dot: 'bg-amber-400',
		text: 'text-amber-400',
		hex: '#fbbf24'
	},
	moved: {
		label: 'moved',
		badge: 'border-violet-500/30 bg-violet-500/10 text-violet-400',
		dot: 'bg-violet-400',
		text: 'text-violet-400',
		hex: '#a78bfa'
	}
};

export const LEGEND_ORDER: Status[] = ['in_sync', 'untracked', 'phantom', 'drift', 'moved'];

export function statusMeta(status: Status): StatusMeta {
	return STATUS_META[status] ?? STATUS_META.in_sync;
}

// statusLabelUpper returns the uppercase status used on node badges ("IN SYNC").
export function statusLabelUpper(status: Status): string {
	return statusMeta(status).label.toUpperCase();
}
