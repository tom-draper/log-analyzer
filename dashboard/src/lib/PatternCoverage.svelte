<script lang="ts">
  type PatternCount = { pattern: string; count: number; unmatched: boolean };

  let { data }: { data: Data } = $props();
  let coverage = $derived.by(() => {
    const counts: Record<string, number> = {};
    let unmatched = 0;
    for (const extraction of data.extraction) {
      if (extraction.pattern === "") unmatched++;
      else counts[extraction.pattern] = (counts[extraction.pattern] ?? 0) + 1;
    }
    const patterns: PatternCount[] = Object.entries(counts)
      .map(([pattern, count]) => ({ pattern, count, unmatched: false }))
      .sort((a, b) => b.count - a.count);
    if (unmatched > 0) patterns.unshift({ pattern: "Unmatched lines", count: unmatched, unmatched: true });
    return { patterns, unmatched, matched: data.extraction.length - unmatched };
  });
  let matchRate = $derived(data.extraction.length === 0 ? 0 : (coverage.matched / data.extraction.length) * 100);
</script>

<section class="coverage">
  <div class="heading"><div><h2>Pattern coverage</h2><span>{coverage.matched.toLocaleString()} of {data.extraction.length.toLocaleString()} lines matched</span></div><strong class:warning={coverage.unmatched > 0}>{matchRate.toFixed(1)}%</strong></div>
  <div class="bar"><span style={`width: ${matchRate}%`}></span></div>
  <div class="patterns">
    {#each coverage.patterns as item}
      <div class:unmatched={item.unmatched} class="pattern"><code title={item.pattern}>{item.pattern}</code><span>{item.count.toLocaleString()} <small>({((item.count / data.extraction.length) * 100).toFixed(1)}%)</small></span></div>
    {/each}
  </div>
</section>

<style scoped>
  .coverage { border: 1px solid var(--border); border-radius: var(--radius); background: var(--card); padding: 1.25rem 1.5rem; margin: 1.5rem 0; } .heading { display: flex; justify-content: space-between; align-items: center; } h2 { margin: 0; font-size: 1rem; } .heading span, small { color: var(--muted); font-size: .8rem; } strong { color: var(--highlight); font-size: 1.4rem; } strong.warning { color: #e0b34b; }
  .bar { height: .5rem; border-radius: 1rem; background: #2b2b2b; overflow: hidden; margin: 1rem 0; } .bar span { display: block; height: 100%; background: var(--highlight); } .patterns { display: grid; gap: .45rem; } .pattern { display: flex; gap: 1rem; justify-content: space-between; font-size: .8rem; min-width: 0; } .pattern code { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #bbb; } .pattern.unmatched code { color: #dd7178; }
</style>
