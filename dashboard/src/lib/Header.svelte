<script lang="ts">
  let { multiTypeTokens, failedLines, lineCount }: {
    multiTypeTokens: DataTypes;
    failedLines: FailedLines;
    lineCount: number;
  } = $props();

  function scrollToBottom() {
    window.scrollTo(0, document.body.scrollHeight);
  }
</script>

<div class="header">
  <div class="title">
    <div class="logo">Log Analyzer</div>
    <div class="line-count">{lineCount.toLocaleString()} lines</div>
  </div>
  <div class="notifications">
    {#if Object.keys(multiTypeTokens).length >= 1}
      <button onclick={scrollToBottom} class="warning"
        >{Object.keys(multiTypeTokens).length}
        {Object.keys(multiTypeTokens).length == 1
          ? "warning"
          : "warnings"}</button
      >
    {/if}
    {#if Object.keys(failedLines).length >= 1}
      <button onclick={scrollToBottom} class="error"
        >{Object.keys(failedLines).length}
        {Object.keys(failedLines).length == 1 ? "error" : "errors"}</button
      >
    {/if}
  </div>
</div>

<style scoped>
  .title {
    margin: 0;
    padding-bottom: .75rem;
  }
  .logo {
    font-size: 1.5rem;
    font-weight: 600;
    letter-spacing: -0.04em;
  }
  .line-count {
    font-size: .875rem;
    color: #a1a1aa;
    margin-top: .2rem;
  }
  .header {
    display: flex;
    align-items: flex-start;
  }

  .notifications {
    margin-left: auto;
    display: flex;
    flex-wrap: wrap;
    justify-content: flex-end;
    gap: .4rem;
  }
  button {
    padding: .35rem .55rem;
    border-radius: 6px;
    height: min-content;
    font-size: 0.9rem;
    outline: none;
  }

  .error {
    background: #2a1518;
    color: #fda4af;
    border: 1px solid #9f1239;
  }
  .warning {
    background: #292319;
    color: #fbbf24;
    border: 1px solid #854d0e;
  }
  .error,
  .warning {
    margin-left: 0;
  }
  @media (max-width: 520px) {
    .header { display: grid; gap: .5rem; }
    .notifications { margin-left: 0; justify-content: flex-start; }
  }
</style>
