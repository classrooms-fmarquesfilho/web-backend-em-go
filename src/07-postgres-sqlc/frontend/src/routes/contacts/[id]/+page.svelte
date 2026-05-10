<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { getContact, deleteContact } from '$lib/api.js';

  /** @type {import('$lib/api.js').Contact | null} */
  let contact = null;
  let loading = true;
  let error = '';

  $: id = Number($page.params.id);

  async function load() {
    loading = true;
    error = '';
    try {
      contact = await getContact(id);
    } catch (e) {
      error = String(e.message ?? e);
      contact = null;
    } finally {
      loading = false;
    }
  }

  async function handleDelete() {
    if (!contact) return;
    if (!confirm(`Apagar contato "${contact.name}"?`)) return;
    try {
      await deleteContact(contact.id);
      goto('/');
    } catch (e) {
      error = String(e.message ?? e);
    }
  }

  onMount(load);
</script>

<a href="/" class="back">← Todos os contatos</a>

{#if loading}
  <div class="empty">Carregando...</div>
{:else if error}
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
</style>
