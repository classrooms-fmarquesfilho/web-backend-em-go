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

/**
 * @typedef {Object} Phone
 * @property {number} id
 * @property {number} contact_id
 * @property {string} label
 * @property {string} number
 * @property {string} created_at
 */

/**
 * @typedef {Object} ContactWithPhones
 * @property {number} id
 * @property {string} name
 * @property {string} email
 * @property {{ id: number, label: string, number: string }[]} phones
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

// ── Contatos ────────────────────────────────────────────────────────────────

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

// ── Telefones ───────────────────────────────────────────────────────────────

/** @param {number} contactId @returns {Promise<Phone[]>} */
export const listPhones = (contactId) =>
  http(`/contacts/${contactId}/phones`);

/**
 * @param {number} contactId
 * @param {{ label: string, number: string }} payload
 * @returns {Promise<Phone>}
 */
export const createPhone = (contactId, payload) =>
  http(`/contacts/${contactId}/phones`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });

/** @param {number} contactId @param {number} phoneId @returns {Promise<null>} */
export const deletePhone = (contactId, phoneId) =>
  http(`/contacts/${contactId}/phones/${phoneId}`, { method: 'DELETE' });

// ── JOIN agregado ───────────────────────────────────────────────────────────

/**
 * Consome o endpoint que executa LEFT JOIN no banco e agrega no Go.
 * Útil para a tela "todos os contatos com seus telefones".
 *
 * @returns {Promise<ContactWithPhones[]>}
 */
export const listContactsWithPhones = () => http('/contacts-with-phones');
