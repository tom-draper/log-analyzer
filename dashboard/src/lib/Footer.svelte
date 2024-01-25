<div class="footer">
  <a
    class="github-logo-container"
    href="https://github.com/tom-draper/log-analyzer"
  >
    <img class="github-icon" src="./github.png" alt="" />
  </a>
  <div class="logo-container">
    <div>Log Analyzer</div>
    <img class="icon" src="./icon.png" alt="" />
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
    /* margin: 0.5em 0; */
  }
  .icon {
    width: 1.4em;
    margin-left: 0.5em;
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
