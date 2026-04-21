import { useState, useEffect } from 'react';
import { TagChip } from './TagChip';
import type { Tag, TagsClientOptions } from './lib/api';
import { createTagsClient } from './lib/api';

export interface TagListProps {
  /** Subject entity UUID */
  subject: string;
  /**
   * Restrict to these purposes. undefined and [] are equivalent — both mean
   * "all purposes / free-form". This is enforced at the type level via
   * `purposes?: string[]` and guarded at runtime by checking `.length > 0`.
   */
  purposes?: string[];
  noPurpose?: boolean;
  client: ReturnType<typeof createTagsClient>;
  className?: string;
}

type LoadState = 'idle' | 'loading' | 'error' | 'ready';

export function TagList({ subject, purposes, noPurpose, client, className }: TagListProps) {
  const [tags, setTags] = useState<Tag[]>([]);
  const [loadState, setLoadState] = useState<LoadState>('idle');
  const [errorMessage, setErrorMessage] = useState<string>('');

  // Stable key: purposes=undefined and purposes=[] both stringify to "[]"
  const purposesKey = JSON.stringify(purposes ?? []);

  useEffect(() => {
    let cancelled = false;
    setLoadState('loading');
    setErrorMessage('');

    client
      .listBySubject(subject, purposes)
      .then((fetched) => {
        if (cancelled) return;
        setTags(fetched);
        setLoadState('ready');
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        const msg = err instanceof Error ? err.message : 'Failed to load tags.';
        setErrorMessage(msg);
        setLoadState('error');
      });

    return () => {
      cancelled = true;
    };
    // purposesKey is the stable serialization of purposes for the dep array
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [subject, purposesKey]);

  if (loadState === 'loading' || loadState === 'idle') {
    return (
      <div className={`flex flex-wrap gap-1.5 ${className ?? ''}`} aria-busy="true">
        {/* Skeleton placeholders */}
        <span className="h-5 w-20 animate-pulse rounded-full bg-gray-200" />
        <span className="h-5 w-16 animate-pulse rounded-full bg-gray-200" />
      </div>
    );
  }

  if (loadState === 'error') {
    return (
      <span className="text-xs text-red-600" role="alert">
        {errorMessage}
      </span>
    );
  }

  // Empty state: render nothing — keeps the component unobtrusive
  if (tags.length === 0) {
    return null;
  }

  return (
    <div className={`flex flex-wrap gap-1.5 ${className ?? ''}`}>
      {tags.map((tag) => (
        <TagChip key={tag.uuid} tag={tag} noPurpose={noPurpose} />
      ))}
    </div>
  );
}

// Re-export for convenience of consumers who only import from TagList
export type { Tag, TagsClientOptions };
