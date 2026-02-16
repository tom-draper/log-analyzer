<script lang="ts">
  import ColorPicker, { ChromeVariant } from "svelte-awesome-color-picker";

  let { multiTypeTokens, failedLines, lineCount, hex = $bindable() }: {
    multiTypeTokens: DataTypes;
    failedLines: FailedLines;
    lineCount: number;
    hex: string;
  } = $props();

  $effect(() => {
    document.documentElement.style.setProperty("--highlight", hex);
  });

  function scrollToBottom() {
    window.scrollTo(0, document.body.scrollHeight);
  }
</script>

<div class="header">
  <div class="title">
    <!-- <img src="./logo.png" style="width: 10em;"/> -->
    <div class="logo-container">
      <div class="logo">Log Analyzer</div>
      <div class="color-picker">
        <ColorPicker
          bind:hex
          components={ChromeVariant}
          sliderDirection="horizontal"
          label=""
        />
      </div>
    </div>
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
    margin: 0 0 20px;
  }
  .logo-container {
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .color-picker {
    --input-size: 14px;
    margin-bottom: 7px;
  }
  .logo {
    font-family: "Poppins";
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
    display: flex;
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
  .error,
  .warning {
    margin-left: 10px;
  }
</style>
