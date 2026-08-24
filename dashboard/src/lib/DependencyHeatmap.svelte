<script lang="ts">
  type AxisValue = { value: string; count: number };

  let { data, token, dependentToken }: { data: Data; token: string; dependentToken: string } = $props();
  let matrix = $derived.by(() => {
    const cells: Record<string, Record<string, number>> = {};
    const totals: Record<string, number> = {};
    const dependentTotals: Record<string, number> = {};
    for (const extraction of data.extraction) {
      const left = extraction.params[token]?.value;
      const top = extraction.params[dependentToken]?.value;
      if (left === undefined || top === undefined) continue;
      const l = String(left), t = String(top);
      cells[l] ??= {}; cells[l][t] = (cells[l][t] ?? 0) + 1;
      totals[l] = (totals[l] ?? 0) + 1; dependentTotals[t] = (dependentTotals[t] ?? 0) + 1;
    }
    const sort = (values: Record<string, number>): AxisValue[] => Object.entries(values).map(([value, count]) => ({ value, count })).sort((a, b) => b.count - a.count).slice(0, 8);
    return { cells, rows: sort(totals), columns: sort(dependentTotals) };
  });
  let maximum = $derived(Math.max(1, ...Object.values(matrix.cells).flatMap((row) => Object.values(row))));
</script>

{#if matrix.rows.length > 0}
  <section class="heatmap">
    <div class="heading">{token} by {dependentToken}</div>
    <div class="scroll"><table><thead><tr><th>{token}</th>{#each matrix.columns as column}<th title={column.value}>{column.value}</th>{/each}</tr></thead>
      <tbody>{#each matrix.rows as row}<tr><th title={row.value}>{row.value}</th>{#each matrix.columns as column}
        {@const count = matrix.cells[row.value]?.[column.value] ?? 0}
        <td style={`--intensity: ${count / maximum}`} title={`${row.value} × ${column.value}: ${count}`}>{count || ""}</td>
      {/each}</tr>{/each}</tbody></table></div>
  </section>
{/if}

<style scoped>
  .heatmap { margin: 1.75rem 0; } .heading { font-size: .9rem; font-weight: 600; margin-bottom: .5rem; } .scroll { overflow: auto; }
  table { border-collapse: separate; border-spacing: 3px; font-size: .75rem; min-width: 100%; } th { color: #aaa; font-weight: 500; max-width: 8rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; } td { background: color-mix(in srgb, var(--highlight) calc(var(--intensity) * 85%), #202020); color: #111; font-weight: 600; min-width: 2rem; height: 1.8rem; text-align: center; border-radius: 2px; }
</style>
