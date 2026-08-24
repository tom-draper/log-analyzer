<script lang="ts">
  type Row = { value: string; count: number };

  function valueRows(data: Data, token: string): Row[] {
    const counts: Record<string, number> = {};
    for (const extraction of data.extraction) {
      const value = extraction.params[token]?.value;
      if (value === undefined) continue;
      const label = String(value);
      counts[label] = (counts[label] ?? 0) + 1;
    }
    return Object.entries(counts)
      .map(([value, count]) => ({ value, count }))
      .sort((a, b) => b.count - a.count || a.value.localeCompare(b.value))
      .slice(0, 10);
  }

  let { data, token }: { data: Data; token: string } = $props();
  let selectedValue = $state<string | null>(null);
  let rows = $derived(valueRows(data, token));
  let examples = $derived(selectedValue === null ? [] : data.extraction.filter((e) => String(e.params[token]?.value) === selectedValue).slice(-5).reverse());
</script>

{#if rows.length > 0}
  <section class="top-values">
    <div class="heading">Top values <span>select a row for recent examples</span></div>
    <table>
      <thead><tr><th>Value</th><th>Count</th><th>Share</th></tr></thead>
      <tbody>
        {#each rows as row}
          <tr class:selected={selectedValue === row.value}>
            <td><button onclick={() => selectedValue = selectedValue === row.value ? null : row.value}>{row.value}</button></td>
            <td>{row.count.toLocaleString()}</td>
            <td>{((row.count / data.extraction.length) * 100).toFixed(1)}%</td>
          </tr>
        {/each}
      </tbody>
    </table>
    {#if selectedValue !== null}
      <div class="examples">
        <div>Recent lines for <code>{selectedValue}</code></div>
        {#each examples as example}
          <code class="line">{example.line}</code>
        {/each}
      </div>
    {/if}
  </section>
{/if}

<style scoped>
  .top-values { margin: 1.75rem 0; overflow: auto; } .heading { font-size: .9rem; font-weight: 600; margin-bottom: .5rem; }
  .heading span { color: #888; font-weight: 400; font-size: .85rem; margin-left: .4rem; } table { width: 100%; border-collapse: collapse; font-size: .85rem; }
  th { color: #888; text-align: right; font-weight: 500; } th:first-child, td:first-child { text-align: left; } td, th { border-bottom: 1px solid #292929; padding: .4rem .25rem; }
  button { background: none; color: var(--highlight); border: 0; padding: 0; font: inherit; text-align: left; max-width: 32rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  tr.selected { background: #ffffff0b; } .examples { background: #171717; border-radius: 4px; margin-top: .75rem; padding: .75rem; font-size: .8rem; }
  .line { display: block; color: #aaa; margin-top: .4rem; overflow-wrap: anywhere; }
</style>
