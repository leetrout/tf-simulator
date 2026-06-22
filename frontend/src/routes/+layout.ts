// Prerender the app to a static bundle (served at "/" by the Go binary, with
// assets under "/static"). SSR must stay on so prerendered HTML references
// assets via the absolute `paths.base` prefix rather than relative paths.
export const prerender = true;
