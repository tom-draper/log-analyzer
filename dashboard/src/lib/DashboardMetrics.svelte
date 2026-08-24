<script lang="ts">
  type Metric = { label: string; value: string; tone?: "neutral" | "warning" | "danger" };

  let { data }: { data: Data } = $props();
  let metrics = $derived.by(() => {
    const total = data.extraction.length;
    const unmatched = data.extraction.filter((item) => item.pattern === "").length;
    const tokenCounts: Record<string, number> = {};
    const times: number[] = [];
    for (const extraction of data.extraction) {
      for (const [token, param] of Object.entries(extraction.params)) {
        tokenCounts[token] = (tokenCounts[token] ?? 0) + 1;
        if (param.type === "time") {
          const time = new Date(param.value as string).getTime();
          if (!Number.isNaN(time)) times.push(time);
        }
      }
    }
    const topToken = Object.entries(tokenCounts).sort((a, b) => b[1] - a[1])[0]?.[0] ?? "—";
    const timeRange = times.length < 2 ? "—" : `${new Date(Math.min(...times)).toLocaleDateString()} – ${new Date(Math.max(...times)).toLocaleDateString()}`;
    const matchRate = total === 0 ? "—" : `${(((total - unmatched) / total) * 100).toFixed(1)}%`;
    return [
      { label: "Lines", value: total.toLocaleString() },
      { label: "Match rate", value: matchRate, tone: unmatched ? "warning" : "neutral" },
      { label: "Unmatched", value: unmatched.toLocaleString(), tone: unmatched ? "danger" : "neutral" },
      { label: "Time range", value: timeRange },
      { label: "Top token", value: topToken },
    ] as Metric[];
  });
</script>

<section class="metrics" aria-label="Dashboard summary">
  {#each metrics as metric}
    <div class={`metric ${metric.tone ?? "neutral"}`}><span>{metric.label}</span><strong title={metric.value}>{metric.value}</strong></div>
  {/each}
</section>

<style scoped>
  .metrics { display: grid; grid-template-columns: repeat(5, minmax(0, 1fr)); gap: .75rem; margin: 0 0 1.25rem; }
  .metric { min-width: 0; border: 1px solid var(--border); border-radius: var(--radius); background: var(--card); padding: .75rem .9rem; display: grid; gap: .2rem; }
  .metric span { color: var(--muted); font-size: .72rem; } .metric strong { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: .92rem; }
  .metric.warning strong { color: #fbbf24; } .metric.danger strong { color: #fda4af; }
  @media (max-width: 900px) { .metrics { grid-template-columns: repeat(3, 1fr); } }
  @media (max-width: 600px) { .metrics { grid-template-columns: repeat(2, 1fr); } }
</style>
