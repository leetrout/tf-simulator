export function graphJSON(data, opts = {}) {
    const errors = [];
    // accept { nodes, edges } / d3-style { nodes, links } / a bare array of nodes
    let rawNodes = [];
    let rawEdges = [];
    if (Array.isArray(data)) {
        rawNodes = data;
    }
    else if (data && typeof data === "object") {
        rawNodes = (data.nodes || data.vertices || []);
        rawEdges = (data.edges || data.links || data.relationships || []);
    }
    if (!Array.isArray(rawNodes)) {
        errors.push("`nodes` is not an array");
        rawNodes = [];
    }
    if (!Array.isArray(rawEdges)) {
        errors.push("`edges` is not an array");
        rawEdges = [];
    }
    const idKey = opts.id || "id";
    const labelKey = "label" in opts ? opts.label : null;
    const srcKey = opts.source || "source";
    const tgtKey = opts.target || "target";
    const srcFieldKey = opts.sourceField || "sourceField";
    const tgtFieldKey = opts.targetField || "targetField";
    const maxLen = opts.maxValueLength == null ? 60 : opts.maxValueLength;
    const clip = (s, n) => {
        const str = String(s);
        return str.length > n ? str.slice(0, n - 1) + "…" : str;
    };
    const fmt = (v) => {
        if (v === null)
            return "null";
        if (v === undefined)
            return "";
        const t = typeof v;
        if (t === "string")
            return clip(v, maxLen);
        if (t === "number" || t === "boolean" || t === "bigint")
            return String(v);
        try {
            return clip(JSON.stringify(v), maxLen);
        }
        catch {
            return clip(String(v), maxLen);
        }
    };
    const getId = typeof idKey === "function"
        ? idKey
        : (n) => (n && n[idKey] != null ? n[idKey] : undefined);
    const getLabel = typeof labelKey === "function"
        ? labelKey
        : (n) => {
            if (labelKey && n && n[labelKey] != null)
                return n[labelKey];
            if (!n)
                return "";
            return n.label != null ? n.label : n.name != null ? n.name : n.title != null ? n.title : getId(n);
        };
    const normCol = (c, node) => ({
        name: String(c.name),
        type: c.type != null ? String(c.type) : node ? fmt(node[c.name]) : "",
        pk: !!c.pk,
        fk: !!c.fk,
        nn: !!c.nn,
        unique: !!c.unique,
    });
    // derive the rows shown inside a node's card
    const fieldsOf = (node) => {
        if (typeof opts.fields === "function") {
            return (opts.fields(node) || []).map((f) => (typeof f === "string" ? specFromKey(f, node) : normCol(f, node)));
        }
        let keys;
        if (Array.isArray(opts.fields)) {
            keys = opts.fields;
        }
        else {
            const skip = new Set(["columns"]);
            if (typeof idKey === "string")
                skip.add(idKey);
            if (typeof labelKey === "string")
                skip.add(labelKey);
            // when no explicit label key, the header is derived from the first of
            // label/name/title present — skip that one so it isn't shown twice
            else if (node) {
                for (const cand of ["label", "name", "title"]) {
                    if (node[cand] != null) {
                        skip.add(cand);
                        break;
                    }
                }
            }
            keys = Object.keys(node || {}).filter((k) => !skip.has(k));
        }
        return keys.map((k) => (typeof k === "string" ? specFromKey(k, node) : normCol(k, node)));
    };
    const specFromKey = (k, node) => ({
        name: String(k),
        type: fmt(node ? node[k] : undefined),
        pk: false,
        fk: false,
        nn: false,
        unique: false,
    });
    const tables = [];
    const seenKeys = new Set();
    for (const node of rawNodes) {
        const idVal = getId(node);
        if (idVal == null) {
            errors.push("node is missing an id; skipped");
            continue;
        }
        const key = String(idVal).toLowerCase();
        if (seenKeys.has(key)) {
            errors.push(`duplicate node id "${idVal}"; later one ignored`);
            continue;
        }
        seenKeys.add(key);
        const declared = node["columns"];
        const columns = Array.isArray(declared)
            ? declared.map((c) => normCol(c, node))
            : fieldsOf(node);
        tables.push({ name: String(getLabel(node)), key, columns, x: NaN, y: NaN, w: 0, h: 0, data: node });
    }
    // read an edge endpoint — supports d3-style object endpoints and from/to aliases
    const endpoint = (edge, key, alias) => {
        let v = edge ? edge[key] : undefined;
        if (v == null && alias)
            v = edge ? edge[alias] : undefined;
        if (v && typeof v === "object")
            return getId(v);
        return v;
    };
    const relations = [];
    for (const edge of rawEdges) {
        const s = endpoint(edge, srcKey, "from");
        const t = endpoint(edge, tgtKey, "to");
        if (s == null || t == null) {
            errors.push("edge is missing source/target; skipped");
            continue;
        }
        if (!seenKeys.has(String(s).toLowerCase()) || !seenKeys.has(String(t).toLowerCase())) {
            errors.push(`edge references unknown node (${s} → ${t}); skipped`);
            continue;
        }
        const sf = edge ? edge[srcFieldKey] : undefined;
        const tf = edge ? edge[tgtFieldKey] : undefined;
        const rel = {
            fromTable: String(s),
            fromCols: sf != null ? [String(sf)] : [],
            toTable: String(t),
            toCols: tf != null ? [String(tf)] : [],
        };
        const label = edge ? (edge["label"] != null ? edge["label"] : edge["type"]) : undefined;
        if (label != null)
            rel.label = String(label);
        relations.push(rel);
    }
    return { tables, relations, errors };
}
//# sourceMappingURL=graph-json.js.map