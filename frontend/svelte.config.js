import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const isLocal = process.env.SK_LOCAL == 'development';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter(),
		paths: {
			// Go serves the built assets under /static
			base: isLocal ? '' : '/static',
			// Emit absolute, base-prefixed asset URLs (e.g. /static/internal/...).
			// The Go server serves index.html at "/" but assets only under
			// "/static", so relative URLs would resolve to the wrong path.
			relative: false
		},
		// SvelteKit uses _app by default but Go (embed.FS) will not serve
		// directories with leading underscores, so use a plain dir name.
		appDir: 'internal'
	}
};

export default config;
