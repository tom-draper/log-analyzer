<script lang="ts">
  let { data }: { data: Data } = $props();
  let query = $state("");
  let failuresOnly = $state(false);
  let limit = $state(50);
  let events = $derived.by(() => data.extraction.filter((event) => {
    if (failuresOnly && event.pattern !== "") return false;
    return query === "" || event.line.toLowerCase().includes(query.toLowerCase());
  }).slice(-limit).reverse());
</script>

<section class="explorer">
  <div class="heading"><h2>Event explorer</h2><span>latest matching lines</span></div>
  <div class="controls"><input bind:value={query} placeholder="Filter raw log lines" aria-label="Filter raw log lines" /><label><input type="checkbox" bind:checked={failuresOnly} /> unmatched only</label></div>
  <div class="events">{#each events as event}
    <div class:failed={event.pattern === ""} class="event"><span>#{event.lineNumber + 1}</span><code>{event.line}</code></div>
  {:else}<p>No matching events.</p>{/each}</div>
  {#if events.length === limit}<button onclick={() => limit += 50}>Show 50 more</button>{/if}
</section>

<style scoped>
  .explorer { border: 1px solid var(--border); border-radius: var(--radius); background: var(--card); padding: 1.25rem 1.5rem; margin: 1.5rem 0; } .heading { display: flex; gap: .6rem; align-items: baseline; } h2 { margin: 0; font-size: 1rem; } .heading span, p { color: var(--muted); font-size: .85rem; }
  .controls { display: flex; gap: 1rem; margin: 1rem 0; align-items: center; font-size: .85rem; } .controls > input { flex: 1; } input { background: #181818; color: var(--text); border: 1px solid #333; border-radius: 4px; padding: .45rem .6rem; }
  .events { max-height: 22rem; overflow: auto; background: #151515; } .event { display: flex; gap: .7rem; border-bottom: 1px solid #252525; padding: .45rem .6rem; font-size: .8rem; } .event span { color: #777; min-width: 3rem; } .event code { overflow-wrap: anywhere; } .failed { border-left: 3px solid #dd7178; } button { margin-top: .75rem; font-size: .8rem; padding: .4rem .6rem; }
  @media (max-width: 600px) { .controls { display: grid; } }
</style>
