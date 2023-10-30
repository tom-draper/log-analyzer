<script lang="ts">
  import { onMount } from "svelte";

  function getFailedExtractions(extraction: Extraction[]): Extraction[] {
    let failedLines = []
      for (let i = 0; i < data.extraction.length; i++) {
        if (data.extraction[i].params.length === 0) {
          failedLines.push(data.extraction[i])
        }
        
      }
      return failedLines
  }

  let failedLines: Extraction[];
  onMount(() => {
    failedLines = getFailedExtractions(data.extraction)
  })
  export let data: Data;
</script>

{#if failedLines !== undefined}
  <div class="card">
    {#each failedLines.slice(0, 100) as failedLine}
      <div class="line-container">
        <div class="lineNumber">{failedLine.lineNumber}</div>
        <div class="line">{failedLine.line}</div>
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

