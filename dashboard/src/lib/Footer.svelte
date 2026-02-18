<script lang="ts">
  import githubIcon from '../assets/github.png'
  export let hex: string;
</script>

<div class="footer">
  <a
    class="github-logo-container"
    href="https://github.com/tom-draper/log-analyzer"
  >
    <img class="github-icon" src={githubIcon} alt="" />
  </a>
  <div class="logo-container">
    <div>Log Analyzer</div>
    <div class="dot" style="background: {hex}"></div>
  </div>
</div>

<style scoped>
  .footer {
    display: grid;
    place-items: center;
    font-size: 0.8em;
    padding: 1em 0 2em;
    font-weight: 600;
  }
  .logo-container {
    display: flex;
    align-items: center;
    gap: 6px;
    font-family: "Poppins";
  }
  .dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    flex-shrink: 0;
    align-self: flex-end;
    margin-bottom: 4px;
  }
  .github-icon:hover {
    filter: contrast(0.3);
  }
  .github-icon {
    width: 2em;
    filter: contrast(0.5);
    padding-bottom: 1.5em;
    transition: filter 0.1s;
  }
</style>
