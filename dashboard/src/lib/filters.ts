export type DashboardFilters = {
  start: string;
  end: string;
  token: string;
  value: string | null;
  source: string | null;
  query: string;
  failuresOnly: boolean;
};

export function emptyFilters(): DashboardFilters {
  return { start: "", end: "", token: "", value: null, source: null, query: "", failuresOnly: false };
}

export function filterRecords(records: Extraction[], filters: DashboardFilters): Extraction[] {
  const start = filters.start ? new Date(filters.start).getTime() : -Infinity;
  const end = filters.end ? new Date(filters.end).getTime() : Infinity;
  const query = filters.query.toLowerCase();
  return records.filter((record) => {
    if (filters.source !== null && (record.sourceFile ?? "") !== filters.source) return false;
    if (filters.failuresOnly && record.pattern !== "") return false;
    if (query && !record.line.toLowerCase().includes(query)) return false;
    if (filters.token) {
      const param = record.params[filters.token];
      if (!param || (filters.value !== null && String(param.value) !== filters.value)) return false;
    }
    if (filters.start || filters.end) {
      const param = Object.values(record.params).find(param => param.type === "time");
      if (!param) return false;
      const timestamp = new Date(param.value).getTime();
      if (!Number.isFinite(timestamp) || !(timestamp >= start && timestamp <= end)) return false;
    }
    return true;
  });
}

export function exportRecords(records: Extraction[]): string {
  return JSON.stringify(records, null, 2);
}
