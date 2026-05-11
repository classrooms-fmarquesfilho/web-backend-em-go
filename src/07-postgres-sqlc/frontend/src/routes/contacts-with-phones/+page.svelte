<script>
  import { onMount } from 'svelte';
  import { listContactsWithPhones } from '$lib/api.js';

  /** @type {import('$lib/api.js').ContactWithPhones[]} */
  let data = [];
  let loading = true;
  let error = '';
  let elapsedMs = 0;
  let showRaw = false;

  async function load() {
    loading = true;
    error = '';
    const t0 = performance.now();
    try {
      const result = await listContactsWithPhones();
      data = result ?? [];
    } catch (e) {
      error = String(e.message ?? e);
    } finally {
      elapsedMs = Math.round(performance.now() - t0);
      loading = false;
    }
  }

  onMount(load);
</script>

<a href="/" class="back">← Lista simples</a>

<header class="hero">
  <h1>Contatos com telefones</h1>
  <p class="tagline">
    Esta tela consome um <strong>único</strong> endpoint
    <code>GET /contacts-with-phones</code> que executa um <code>LEFT JOIN</code>
    no banco e devolve a estrutura aninhada já pronta.
  </p>
</header>

{#if loading}
  <div class="empty">Carregando...</div>
{:else if error}
  <div class="error">⚠️ {error}</div>
{:else}
  <div class="stats">
    <span class="stat">
      <strong>{data.length}</strong> contatos
    </span>
    <span class="stat">
      <strong>{data.reduce((acc, c) => acc + c.phones.length, 0)}</strong> telefones no total
    </span>
    <span class="stat">
      <strong>1</strong> requisição HTTP · <strong>1</strong> query SQL
    </span>
    <span class="stat muted">
      ({elapsedMs}ms)
    </span>
  </div>

  {#if data.length === 0}
    <div class="empty">Nenhum contato cadastrado ainda. <a href="/">Adicionar</a>.</div>
  {:else}
    <div class="grid">
      {#each data as c (c.id)}
        <article class="card">
          <header>
            <h2><a href={`/contacts/${c.id}`}>{c.name}</a></h2>
            <span class="email">{c.email}</span>
          </header>

          {#if c.phones.length === 0}
            <p class="no-phones">Sem telefones</p>
          {:else}
            <ul>
              {#each c.phones as p (p.id)}
                <li>
                  <span class="label">{p.label}</span>
                  <span class="number">{p.number}</span>
                </li>
              {/each}
            </ul>
          {/if}
        </article>
      {/each}
    </div>
  {/if}

  <!-- Aba opcional para mostrar o JSON cru — útil para aula -->
  <section class="debug">
    <button class="toggle" on:click={() => (showRaw = !showRaw)}>
      {showRaw ? '▼' : '▶'} JSON recebido do backend
    </button>
    {#if showRaw}
      <pre>{JSON.stringify(data, null, 2)}</pre>
    {/if}
  </section>
{/if}

<style>
  .back {
    display: inline-block;
    text-decoration: none;
    margin-bottom: 1rem;
    color: var(--muted);
    font-size: 0.9rem;
  }
  .back:hover { color: var(--accent); }

  .hero {
    margin-bottom: 1.2rem;
  }
  .hero h1 {
    margin: 0 0 0.3rem 0;
    font-size: 1.6rem;
  }
  .tagline {
    margin: 0;
    color: var(--muted);
    line-height: 1.5;
  }
  .tagline code {
    background: rgba(184, 92, 56, 0.08);
    color: var(--accent);
    padding: 0.1em 0.35em;
    border-radius: 3px;
    font-size: 0.9em;
  }

  .stats {
    display: flex;
    flex-wrap: wrap;
    gap: 1.2rem;
    padding: 0.7rem 1rem;
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 6px;
    margin-bottom: 1.5rem;
    font-size: 0.92rem;
  }
  .stat strong { color: var(--accent); }
  .stat.muted { margin-left: auto; color: var(--muted); }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 0.8rem;
  }
  .card {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 1rem 1.1rem;
    box-shadow: var(--shadow);
  }
  .card header {
    margin-bottom: 0.7rem;
  }
  .card h2 {
    margin: 0;
    font-size: 1.05rem;
  }
  .card h2 a {
    color: inherit;
    text-decoration: none;
  }
  .card h2 a:hover { color: var(--accent); }
  .email {
    display: block;
    color: var(--muted);
    font-size: 0.85rem;
    margin-top: 0.15rem;
  }
  .no-phones {
    margin: 0;
    padding: 0.3rem 0;
    color: var(--muted);
    font-size: 0.85rem;
    font-style: italic;
  }
  ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  li {
    display: grid;
    grid-template-columns: 5rem 1fr;
    gap: 0.5rem;
    padding: 0.3rem 0;
    border-top: 1px solid var(--border);
  }
  li:first-child { border-top: none; }
  .label {
    font-family: ui-monospace, "Consolas", monospace;
    font-size: 0.78rem;
    color: var(--accent);
    text-transform: lowercase;
    align-self: center;
  }
  .number {
    font-variant-numeric: tabular-nums;
    font-size: 0.92rem;
  }

  .debug {
    margin-top: 2rem;
    border-top: 1px solid var(--border);
    padding-top: 1rem;
  }
  .toggle {
    background: transparent;
    border: none;
    color: var(--muted);
    cursor: pointer;
    padding: 0.3rem 0;
    font-family: ui-monospace, "Consolas", monospace;
    font-size: 0.85rem;
  }
  .toggle:hover { color: var(--accent); }
  pre {
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 1rem;
    overflow-x: auto;
    font-size: 0.78rem;
    line-height: 1.45;
  }
</style>
