// $lib/api.js — Cliente da API Go.
//
// Em desenvolvimento, /api/* é proxiado para http://localhost:8080 (vite.config.js).
// Em produção, ajuste BASE para apontar para a sua API.

const BASE = '/api';

/**
 * @typedef {Object} Contact
 * @property {number} id
 * @property {string} name
 * @property {string} email
 * @property {string} created_at
 */

async function http(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...(options.headers ?? {}) },
    ...options
  });

  // 204 No Content
  if (res.status === 204) return null;

  // Erro: o backend usa Problem Details (RFC 7807)
  if (!res.ok) {
    let detail = res.statusText;
    try {
      const body = await res.json();
      detail = body.detail || body.title || detail;
    } catch {
      /* corpo não-JSON */
    }
    throw new Error(`${res.status}: ${detail}`);
  }

  return res.json();
}

/** @returns {Promise<Contact[]>} */
export const listContacts = () => http('/contacts');

/** @param {number} id @returns {Promise<Contact>} */
export const getContact = (id) => http(`/contacts/${id}`);

/**
 * @param {{ name: string, email: string }} payload
 * @returns {Promise<Contact>}
 */
export const createContact = (payload) =>
  http('/contacts', { method: 'POST', body: JSON.stringify(payload) });

/** @param {number} id @returns {Promise<null>} */
export const deleteContact = (id) =>
  http(`/contacts/${id}`, { method: 'DELETE' });
