<script lang="ts">
  import githubIcon from '../assets/github.png'
</script>

<div class="footer">
  <a
    class="github-logo-container"
    href="https://github.com/tom-draper/log-analyzer"
  >
    <img class="github-icon" src={githubIcon} alt="" />
  </a>
  <div class="logo-container">Log Analyzer</div>
</div>

<style scoped>
  .footer {
    display: grid;
    place-items: center;
    font-size: 0.8em;
    padding: 1rem 0 1.5rem;
    font-weight: 500;
  }
  .logo-container {
    display: flex;
    align-items: center;
    color: #a1a1aa;
  }
  .github-icon:hover {
    filter: contrast(0.3);
  }
  .github-icon {
    width: 2em;
    filter: contrast(0.5);
    padding-bottom: 1rem;
    transition: filter 0.1s;
  }
</style>
