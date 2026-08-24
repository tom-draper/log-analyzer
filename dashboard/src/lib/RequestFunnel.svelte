<script lang="ts">
  type Funnel = {
    total: number;
    clientErrors: number;
    serverErrors: number;
    statusToken: string | null;
  };

  function findStatusToken(data: Data): string | null {
    const counts: Record<string, number> = {};
    for (const extraction of data.extraction) {
      for (const [token, param] of Object.entries(extraction.params)) {
        const value = param.value;
        if (typeof value !== "number" || value < 100 || value > 599) continue;
        if (!/(status|code)/i.test(token)) continue;
        counts[token] = (counts[token] ?? 0) + 1;
      }
    }
    return Object.entries(counts).sort((a, b) => b[1] - a[1])[0]?.[0] ?? null;
  }

  function buildFunnel(data: Data): Funnel {
    const statusToken = findStatusToken(data);
    if (!statusToken) return { total: 0, clientErrors: 0, serverErrors: 0, statusToken: null };

    let total = 0;
    let clientErrors = 0;
    let serverErrors = 0;
    for (const extraction of data.extraction) {
      const value = extraction.params[statusToken]?.value;
      if (typeof value !== "number") continue;
      total++;
      if (value >= 400 && value < 500) clientErrors++;
      if (value >= 500 && value < 600) serverErrors++;
    }
    return { total, clientErrors, serverErrors, statusToken };
  }

  let { data }: { data: Data } = $props();
  let funnel = $derived(buildFunnel(data));
  const percentage = (value: number) => funnel.total === 0 ? "0" : ((value / funnel.total) * 100).toFixed(1);
</script>

{#if funnel.statusToken}
  <section class="funnel" aria-label="Request health funnel">
    <div class="heading">
      <h2>Request health</h2>
      <span>from <code>{funnel.statusToken}</code></span>
    </div>
    <div class="steps">
      <div class="step total">
        <span class="label">Requests</span>
        <strong>{funnel.total.toLocaleString()}</strong>
      </div>
      <div class="step client">
        <span class="label">4xx responses</span>
        <strong>{funnel.clientErrors.toLocaleString()}</strong>
        <small>{percentage(funnel.clientErrors)}%</small>
      </div>
      <div class="step server">
        <span class="label">5xx responses</span>
        <strong>{funnel.serverErrors.toLocaleString()}</strong>
        <small>{percentage(funnel.serverErrors)}%</small>
      </div>
    </div>
  </section>
{/if}

<style scoped>
  .funnel { border: 1px solid var(--border); border-radius: var(--radius); background: var(--card); margin: 0 0 1.5rem; padding: 1.25rem 1.5rem; }
  .heading { display: flex; align-items: baseline; gap: .6rem; margin-bottom: 1rem; }
  h2 { font-size: 1rem; margin: 0; } .heading span, small { color: #888; font-size: .8rem; }
  .steps { display: grid; grid-template-columns: repeat(3, 1fr); gap: .8rem; }
  .step { border-left: 3px solid var(--highlight); padding-left: .8rem; display: grid; }
  .client { border-color: #e0b34b; } .server { border-color: #dd7178; }
  .label { color: #aaa; font-size: .85rem; } strong { font-size: 1.4rem; }
  @media (max-width: 600px) { .steps { grid-template-columns: 1fr; } }
</style>
