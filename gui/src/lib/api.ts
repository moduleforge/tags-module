// API client for the tags service. Zero runtime deps beyond fetch.
//
// Field naming: the wire format uses camelCase (ownerUuid, subjectUuid, createdAt, etc.)
// as confirmed by tags.go tagResponse struct tags. TypeScript types match wire exactly.

export interface Tag {
  uuid: string;
  ownerUuid: string;
  subjectUuid: string;
  purpose: string;
  value: string;
  color?: string;
  createdAt: string;
  updatedAt: string;
}

export interface TagsClientOptions {
  /** Base URL for the API (e.g. "/v1"). Required. */
  baseUrl: string;
  /** Fetch implementation; defaults to global fetch. */
  fetchImpl?: typeof fetch;
  /**
   * Called before every request; return an object of headers to merge.
   * Typically used by consumers to inject Authorization: Bearer <token>.
   */
  headers?: () => Record<string, string> | Promise<Record<string, string>>;
}

interface ErrorResponse {
  error?: string;
  message?: string;
}

async function handleResponse<T>(res: Response): Promise<T> {
  if (res.ok) {
    // 204 No Content has no body
    if (res.status === 204) {
      return undefined as T;
    }
    return res.json() as Promise<T>;
  }

  // Attempt to parse error body as JSON { error: "..." } or { message: "..." }
  let message = res.statusText;
  try {
    const body = (await res.json()) as ErrorResponse;
    message = body.error ?? body.message ?? message;
  } catch {
    // Body wasn't JSON; use status text
  }
  throw new Error(message);
}

export function createTagsClient(opts: TagsClientOptions): {
  listBySubject(subjectUuid: string, purposes?: string[]): Promise<Tag[]>;
  create(input: { subject: string; purpose: string; value: string; color?: string }): Promise<Tag>;
  updateColor(uuid: string, color: string | null): Promise<Tag>;
  remove(uuid: string): Promise<void>;
} {
  const fetchFn = opts.fetchImpl ?? globalThis.fetch;

  async function buildHeaders(): Promise<Record<string, string>> {
    const extra = opts.headers ? await opts.headers() : {};
    return { 'Content-Type': 'application/json', ...extra };
  }

  return {
    async listBySubject(subjectUuid: string, purposes?: string[]): Promise<Tag[]> {
      const headers = await buildHeaders();
      const url = `${opts.baseUrl}/entities/${subjectUuid}/tags`;
      const res = await fetchFn(url, { headers });
      // Server returns { tags: Tag[] }
      const body = await handleResponse<{ tags: Tag[] }>(res);
      const tags = body.tags ?? [];

      // Client-side filter for robustness when multiple purposes are requested.
      // purposes=undefined or purposes=[] means "all" — no filter applied.
      const effectivePurposes = purposes && purposes.length > 0 ? purposes : null;
      if (effectivePurposes === null || effectivePurposes.length === 1) {
        // Server already filtered for 0 or 1 purpose; return as-is.
        return tags;
      }
      // Multiple purposes: filter client-side
      const purposeSet = new Set(effectivePurposes);
      return tags.filter((t) => purposeSet.has(t.purpose));
    },

    async create(input: {
      subject: string;
      purpose: string;
      value: string;
      color?: string;
    }): Promise<Tag> {
      const headers = await buildHeaders();
      const res = await fetchFn(`${opts.baseUrl}/tags`, {
        method: 'POST',
        headers,
        body: JSON.stringify(input),
      });
      return handleResponse<Tag>(res);
    },

    async updateColor(uuid: string, color: string | null): Promise<Tag> {
      const headers = await buildHeaders();
      const res = await fetchFn(`${opts.baseUrl}/tags/${uuid}`, {
        method: 'PUT',
        headers,
        body: JSON.stringify({ color }),
      });
      return handleResponse<Tag>(res);
    },

    async remove(uuid: string): Promise<void> {
      const headers = await buildHeaders();
      const res = await fetchFn(`${opts.baseUrl}/tags/${uuid}`, {
        method: 'DELETE',
        headers,
      });
      return handleResponse<void>(res);
    },
  };
}
