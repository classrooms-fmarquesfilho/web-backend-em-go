<script>
  import { onMount } from 'svelte';
  import { listContacts, createContact, deleteContact } from '$lib/api.js';

  /** @type {import('$lib/api.js').Contact[]} */
  let contacts = [];
  let loading = true;
  let error = '';

  // form
  let name = '';
  let email = '';
  let submitting = false;

  async function load() {
    loading = true;
    error = '';
    try {
      const data = await listContacts();
      contacts = data ?? [];
    } catch (e) {
      error = String(e.message ?? e);
    } finally {
      loading = false;
    }
  }

  async function handleCreate() {
    if (!name.trim() || !email.trim()) return;
    submitting = true;
    error = '';
    try {
      const created = await createContact({ name: name.trim(), email: email.trim() });
      contacts = [created, ...contacts];
      name = '';
      email = '';
    } catch (e) {
      error = String(e.message ?? e);
    } finally {
      submitting = false;
    }
  }

  async function handleDelete(/** @type {number} */ id, /** @type {string} */ contactName) {
    if (!confirm(`Apagar contato "${contactName}"?`)) return;
    error = '';
    try {
      await deleteContact(id);
      contacts = contacts.filter((c) => c.id !== id);
    } catch (e) {
      error = String(e.message ?? e);
    }
  }

  function fmtDate(/** @type {string} */ iso) {
    try {
      return new Date(iso).toLocaleDateString('pt-BR', {
        day: '2-digit',
        month: 'short',
        year: 'numeric'
      });
    } catch {
      return iso;
    }
  }

  onMount(load);
</script>

<section class="form-card">
  <h2>Adicionar contato</h2>

  <div class="row">
    <input
      type="text"
      placeholder="Nome"
      bind:value={name}
      on:keydown={(e) => e.key === 'Enter' && handleCreate()}
    />
    <input
      type="email"
      placeholder="email@exemplo.com"
      bind:value={email}
      on:keydown={(e) => e.key === 'Enter' && handleCreate()}
    />
    <button
      class="primary"
      on:click={handleCreate}
      disabled={submitting || !name.trim() || !email.trim()}
    >
      {submitting ? 'Salvando...' : 'Adicionar'}
    </button>
  </div>
</section>

{#if error}
  <div class="error">⚠️ {error}</div>
{/if}

<section>
  <div class="list-header">
    <h2>Contatos {#if !loading}<span class="count">({contacts.length})</span>{/if}</h2>
    <div class="actions-header">
      <a href="/contacts-with-phones" class="link-aggregate">Ver com telefones →</a>
      <button on:click={load} disabled={loading} title="Recarregar">↻</button>
    </div>
  </div>

  {#if loading}
    <div class="empty">Carregando...</div>
  {:else if contacts.length === 0}
    <div class="empty">
      Nenhum contato ainda.<br /><small>Use o formulário acima para começar.</small>
    </div>
  {:else}
    <ul>
      {#each contacts as c (c.id)}
        <li>
          <a href={`/contacts/${c.id}`} class="info">
            <strong>{c.name}</strong>
            <span class="email">{c.email}</span>
            <span class="date">{fmtDate(c.created_at)}</span>
          </a>
          <button class="danger" on:click={() => handleDelete(c.id, c.name)} title="Apagar">
            ✕
          </button>
        </li>
      {/each}
    </ul>
  {/if}
</section>

<style>
  .form-card {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 1rem 1.1rem;
    margin-bottom: 1.5rem;
    box-shadow: var(--shadow);
  }
  .form-card h2 {
    margin: 0 0 0.7rem 0;
    font-size: 1rem;
    color: var(--muted);
    font-weight: 500;
  }
  .row {
    display: grid;
    grid-template-columns: 1fr 1fr auto;
    gap: 0.5rem;
  }
  @media (max-width: 600px) {
    .row { grid-template-columns: 1fr; }
  }
  .list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.7rem;
  }
  .list-header h2 {
    margin: 0;
    font-size: 1.05rem;
  }
  .count {
    color: var(--muted);
    font-weight: normal;
  }
  .actions-header {
    display: flex;
    align-items: center;
    gap: 0.6rem;
  }
  .link-aggregate {
    color: var(--accent);
    text-decoration: none;
    font-size: 0.88rem;
  }
  .link-aggregate:hover { text-decoration: underline; }
  ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  li {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.7rem 0.9rem;
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 6px;
    margin-bottom: 0.4rem;
    box-shadow: var(--shadow);
  }
  .info {
    flex: 1;
    display: grid;
    grid-template-columns: 1fr auto;
    grid-template-rows: auto auto;
    gap: 0.05rem 0.6rem;
    text-decoration: none;
    color: inherit;
  }
  .info strong { grid-row: 1; grid-column: 1; }
  .email {
    grid-row: 2;
    grid-column: 1;
    color: var(--muted);
    font-size: 0.88rem;
  }
  .date {
    grid-row: 1 / 3;
    grid-column: 2;
    align-self: center;
    color: var(--muted);
    font-size: 0.8rem;
  }
  .info:hover strong { color: var(--accent); }
</style>
