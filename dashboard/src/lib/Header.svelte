<script lang="ts">
  function scrollToBottom() {
    window.scrollTo(0, document.body.scrollHeight);
  }
  export let multiTypeTokens: DataTypes, failedLines: FailedLines, lineCount: number;
</script>

<div class="header">
  <div class="title">
    <!-- <img src="./logo.png" style="width: 10em;"/> -->
    <div class="logo-container">
      <div class="logo">Log Analyzer</div>
      <!-- <img alt="" src="./icon.png" /> -->
    </div>
    <div class="line-count">{lineCount.toLocaleString()} lines</div>
  </div>
  <div class="notifications">
    {#if Object.keys(multiTypeTokens).length >= 1}
    <button on:click={scrollToBottom} class="warning"
      >{Object.keys(multiTypeTokens).length}
      {Object.keys(multiTypeTokens).length == 1
      ? "warning"
      : "warnings"}</button
    >
    {/if}
    {#if Object.keys(failedLines).length >= 1}
    <button on:click={scrollToBottom} class="error"
      >{Object.keys(failedLines).length}
      {Object.keys(failedLines).length == 1
      ? "error"
      : "errors"}</button
    >
    {/if}
  </div>
</div>

<style>
  .title {
    margin: 0 0 20px;
    font-family: Poppins;
  }
  .logo-container {
    display: flex;
    align-items: center;
  }
  /* .logo-container>img {
    margin-left: 0.5em;
    width: 1.25em;
  } */
  .logo {
    font-weight: 600;
  }
  .line-count {
    font-size: 1.4rem;
    color: #888;
    margin-top: 0.5em;
  }
  .header {
    font-size: 2em;
    font-weight: 500;
    display: flex;
  }

  .notifications {
    margin-left: auto;
  }
  button {
    padding: 5px 10px;
    border-radius: 3px;
    height: min-content;
    font-size: 0.9rem;
    outline: none;
  }

  .error {
    background: #271515;
    color: #dd7178;
    border: 1px solid #dd71787d;
  }
  .warning {
    background: #5e4c1589;
    color: #ddd871;
    border: 1px solid #ddd8717d;
  }
</style>

