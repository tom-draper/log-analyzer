<script lang="ts">
  const DISPLAY_LIMIT = 100;
  import Panel from "./Panel.svelte";
  export let failedLines: FailedLines;
</script>

{#if Object.keys(failedLines).length > 0}
  {@const total = Object.keys(failedLines).length}
  <Panel class="card">
    <div class="summary">
      {#if total > DISPLAY_LIMIT}
        Showing {DISPLAY_LIMIT.toLocaleString()} of {total.toLocaleString()} failed lines
      {:else}
        {total.toLocaleString()} failed {total === 1 ? "line" : "lines"}
      {/if}
    </div>
    {#each Object.entries(failedLines).slice(0, DISPLAY_LIMIT) as [lineNumber, line]}
      <div class="line-container">
        <div class="lineNumber">{lineNumber}</div>
        <div class="line">{line}</div>
      </div>
    {/each}
  </Panel>
{/if}

<style scoped>
  .summary {
    font-size: 0.9em;
    color: var(--muted);
    margin-bottom: 1rem;
  }
  .line-container {
    display: flex;
    background: #271515;
    color: #dd7178;
    border-radius: calc(var(--radius) - 2px);
    margin: .35rem 0;
    font-size: 0.9em;
  }
  .lineNumber {
    margin: 6px 0;
    padding: 4px .75rem;
    border-right: 1px solid #5e1e2e;
  }
  .line {
    margin: .55rem .75rem;
  }
</style>
