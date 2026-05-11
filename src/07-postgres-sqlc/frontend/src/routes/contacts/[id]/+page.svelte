<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import {
    getContact,
    deleteContact,
    listPhones,
    createPhone,
    deletePhone
  } from '$lib/api.js';

  /** @type {import('$lib/api.js').Contact | null} */
  let contact = null;
  /** @type {import('$lib/api.js').Phone[]} */
  let phones = [];

  let loading = true;
  let error = '';

  // form de novo telefone
  let newLabel = 'casa';
  let newNumber = '';
  let addingPhone = false;

  $: id = Number($page.params.id);

  async function load() {
    loading = true;
    error = '';
    try {
      // Estratégia simples (1 + 1 = 2 queries). Para listagens grandes
      // de contatos, viraria N+1 — aí entra o JOIN em /contacts-with-phones.
      const [c, p] = await Promise.all([getContact(id), listPhones(id)]);
      contact = c;
      phones = p ?? [];
    } catch (e) {
      error = String(e.message ?? e);
      contact = null;
    } finally {
      loading = false;
    }
  }

  async function handleDelete() {
    if (!contact) return;
    if (!confirm(`Apagar contato "${contact.name}"? Os telefones também serão removidos.`)) return;
    try {
      await deleteContact(contact.id);
      goto('/');
    } catch (e) {
      error = String(e.message ?? e);
    }
  }

  async function handleAddPhone() {
    if (!newNumber.trim()) return;
    addingPhone = true;
    error = '';
    try {
      const created = await createPhone(id, {
        label: newLabel.trim(),
        number: newNumber.trim()
      });
      phones = [...phones, created];
      newLabel = 'casa';
      newNumber = '';
    } catch (e) {
      error = String(e.message ?? e);
    } finally {
      addingPhone = false;
    }
  }

  async function handleDeletePhone(/** @type {number} */ phoneId) {
    if (!confirm('Apagar este telefone?')) return;
    try {
      await deletePhone(id, phoneId);
      phones = phones.filter((p) => p.id !== phoneId);
    } catch (e) {
      error = String(e.message ?? e);
    }
  }

  onMount(load);
</script>

<a href="/" class="back">← Todos os contatos</a>

{#if loading}
  <div class="empty">Carregando...</div>
{:else if error && !contact}
  <div class="error">⚠️ {error}</div>
  <p><a href="/">Voltar à lista</a></p>
{:else if contact}
  <article class="card">
    <h1>{contact.name}</h1>
    <dl>
      <dt>ID</dt>
      <dd>#{contact.id}</dd>

      <dt>Email</dt>
      <dd><a href={`mailto:${contact.email}`}>{contact.email}</a></dd>

      <dt>Criado em</dt>
      <dd>{new Date(contact.created_at).toLocaleString('pt-BR')}</dd>
    </dl>

    <div class="actions">
      <button class="danger" on:click={handleDelete}>Apagar contato</button>
    </div>
  </article>

  <!-- ── Telefones ───────────────────────────────────────────────────── -->
  <section class="phones-section">
    <header class="phones-header">
      <h2>Telefones <span class="count">({phones.length})</span></h2>
    </header>

    {#if error && contact}
      <div class="error">⚠️ {error}</div>
    {/if}

    <div class="add-phone">
      <select bind:value={newLabel} disabled={addingPhone}>
        <option value="casa">casa</option>
        <option value="celular">celular</option>
        <option value="trabalho">trabalho</option>
        <option value="recado">recado</option>
      </select>
      <input
        type="text"
        placeholder="+55 84 99999-9999"
        bind:value={newNumber}
        on:keydown={(e) => e.key === 'Enter' && handleAddPhone()}
        disabled={addingPhone}
      />
      <button
        class="primary"
        on:click={handleAddPhone}
        disabled={addingPhone || !newNumber.trim()}
      >
        {addingPhone ? 'Salvando...' : 'Adicionar'}
      </button>
    </div>

    {#if phones.length === 0}
      <p class="empty-phones">Nenhum telefone cadastrado.</p>
    {:else}
      <ul>
        {#each phones as p (p.id)}
          <li>
            <span class="label">{p.label}</span>
            <span class="number">{p.number}</span>
            <button
              class="icon-danger"
              on:click={() => handleDeletePhone(p.id)}
              title="Apagar"
            >
              ✕
            </button>
          </li>
        {/each}
      </ul>
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

  .card {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: var(--shadow);
    margin-bottom: 1.5rem;
  }
  .card h1 {
    margin: 0 0 1.2rem 0;
    font-size: 1.6rem;
  }
  dl {
    display: grid;
    grid-template-columns: 7rem 1fr;
    gap: 0.4rem 1rem;
    margin: 0 0 1.5rem 0;
  }
  dt { color: var(--muted); font-size: 0.9rem; }
  dd { margin: 0; }

  .actions {
    border-top: 1px solid var(--border);
    padding-top: 1rem;
  }

  /* Phones */
  .phones-section {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 1.2rem 1.5rem;
    box-shadow: var(--shadow);
  }
  .phones-header {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    margin-bottom: 1rem;
  }
  .phones-header h2 {
    margin: 0;
    font-size: 1.15rem;
  }
  .count {
    color: var(--muted);
    font-weight: normal;
    font-size: 0.95rem;
  }
  .add-phone {
    display: grid;
    grid-template-columns: 7rem 1fr auto;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  @media (max-width: 600px) {
    .add-phone { grid-template-columns: 1fr; }
  }
  select {
    padding: 0.55rem 0.7rem;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--card);
    color: inherit;
  }
  .empty-phones {
    text-align: center;
    color: var(--muted);
    padding: 1rem;
    margin: 0;
  }
  ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  li {
    display: grid;
    grid-template-columns: 6rem 1fr auto;
    align-items: center;
    gap: 0.6rem;
    padding: 0.5rem 0.7rem;
    border: 1px solid var(--border);
    border-radius: 6px;
    margin-bottom: 0.35rem;
  }
  .label {
    font-family: ui-monospace, "Consolas", monospace;
    font-size: 0.8rem;
    color: var(--accent);
    text-transform: lowercase;
  }
  .number {
    font-variant-numeric: tabular-nums;
  }
  .icon-danger {
    border: none;
    background: transparent;
    color: var(--danger);
    font-size: 1rem;
    cursor: pointer;
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
  }
  .icon-danger:hover {
    background: rgba(192, 57, 43, 0.08);
  }
</style>
