<script>
  import './styles.css';
</script>

<div class="app">
  <header>
    <a href="/" class="brand">📇 Contatos</a>
    <span class="hint">Lista 04 — DIM0547</span>
  </header>

  <main>
    <slot />
  </main>

  <footer>
    <small>Webapp de demonstração da API Go — <a href="https://github.com/webii-2026-1" target="_blank" rel="noopener">UFRN/DIMAp</a></small>
  </footer>
</div>

<style>
  .app {
    max-width: 720px;
    margin: 0 auto;
    padding: 1.5rem 1rem 4rem;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }
  header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    border-bottom: 1px solid var(--border);
    padding-bottom: 0.75rem;
    margin-bottom: 1.5rem;
  }
  .brand {
    font-size: 1.4rem;
    font-weight: 600;
    color: var(--fg);
    text-decoration: none;
  }
  .hint {
    font-size: 0.85rem;
    color: var(--muted);
  }
  main {
    flex: 1;
  }
  footer {
    margin-top: 3rem;
    color: var(--muted);
  }
</style>
