<script lang="ts">
  export let data: Data;
</script>

<div class="card">
  {#each Object.keys(data.extraction.failed).slice(0, 100) as lineNumber}
    <div class="line-container">
      <div class="lineNumber">{lineNumber}</div>
      <div class="line">{data.extraction.failed[lineNumber]}</div>
    </div>
  {/each}
</div>

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

