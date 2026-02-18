<script lang="ts">
  const DISPLAY_LIMIT = 100;
  export let failedLines: FailedLines;
</script>

{#if Object.keys(failedLines).length > 0}
  {@const total = Object.keys(failedLines).length}
  <div class="card">
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
  </div>
{/if}

<style scoped>
  .card {
    border: 1px solid #ffffff24;
    border-radius: 5px;
    margin: 3em 0;
    padding: 2rem;
  }
  .summary {
    font-size: 0.9em;
    color: #888;
    margin-bottom: 1rem;
  }
  .line-container {
    display: flex;
    background: #271515;
    color: #dd7178;
    border-radius: 5px;
    margin: 5px;
    font-size: 0.9em;
  }
  .lineNumber {
    margin: 6px 0;
    padding: 4px 20px;
    border-right: 1px solid #5e1e2e;
  }
  .line {
    margin: 10px 20px;
  }
</style>
