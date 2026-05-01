import type { Story } from '@ladle/react';
import { TagList, type TagListProps } from './TagList';
import type { Tag } from './lib/api';
import { createMockTagsClient } from './lib/mockClient';

const SUBJECT = 'subject-1';

function tag(partial: Partial<Tag> & Pick<Tag, 'uuid' | 'purpose' | 'value'>): Tag {
  return {
    ownerUuid: 'mock-owner',
    subjectUuid: SUBJECT,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
    ...partial,
  };
}

type StoryArgs = Omit<TagListProps, 'client'> & {
  seed: Tag[];
  shouldFail?: boolean;
};

function Render({ seed, shouldFail, ...rest }: StoryArgs) {
  const client = createMockTagsClient({
    initial: seed,
    failOn: shouldFail ? { list: true } : undefined,
  });
  return <TagList {...rest} client={client} />;
}

export const Empty: Story<StoryArgs> = (args) => <Render {...args} />;
Empty.args = { subject: SUBJECT, seed: [] };

export const Few: Story<StoryArgs> = (args) => <Render {...args} />;
Few.args = {
  subject: SUBJECT,
  seed: [
    tag({ uuid: 't1', purpose: 'env', value: 'production', color: '#dc2626' }),
    tag({ uuid: 't2', purpose: 'team', value: 'platform', color: '#2563eb' }),
  ],
};

export const Many: Story<StoryArgs> = (args) => <Render {...args} />;
Many.args = {
  subject: SUBJECT,
  seed: [
    tag({ uuid: 't1', purpose: 'env', value: 'production', color: '#dc2626' }),
    tag({ uuid: 't2', purpose: 'env', value: 'staging', color: '#059669' }),
    tag({ uuid: 't3', purpose: 'env', value: 'dev' }),
    tag({ uuid: 't4', purpose: 'team', value: 'platform', color: '#2563eb' }),
    tag({ uuid: 't5', purpose: 'team', value: 'growth', color: '#7c3aed' }),
    tag({ uuid: 't6', purpose: 'priority', value: 'p0', color: '#f97316' }),
    tag({ uuid: 't7', purpose: 'priority', value: 'p1' }),
    tag({ uuid: 't8', purpose: 'region', value: 'us-east', color: '#fde68a' }),
  ],
};

export const FilteredByPurpose: Story<StoryArgs> = (args) => <Render {...args} />;
FilteredByPurpose.args = {
  subject: SUBJECT,
  purposes: ['env'],
  seed: [
    tag({ uuid: 't1', purpose: 'env', value: 'production', color: '#dc2626' }),
    tag({ uuid: 't2', purpose: 'team', value: 'platform', color: '#2563eb' }),
    tag({ uuid: 't3', purpose: 'env', value: 'staging', color: '#059669' }),
  ],
};

export const NoPurposeDisplay: Story<StoryArgs> = (args) => <Render {...args} />;
NoPurposeDisplay.args = {
  subject: SUBJECT,
  noPurpose: true,
  seed: [
    tag({ uuid: 't1', purpose: 'label', value: 'beta' }),
    tag({ uuid: 't2', purpose: 'label', value: 'alpha' }),
  ],
};

export const LoadError: Story<StoryArgs> = (args) => <Render {...args} />;
LoadError.args = { subject: SUBJECT, seed: [], shouldFail: true };
