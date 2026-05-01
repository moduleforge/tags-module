import type { Story } from '@ladle/react';
import { TagEditor, type TagEditorProps } from './TagEditor';
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

type StoryArgs = Omit<TagEditorProps, 'client' | 'onChange'> & {
  seed: Tag[];
};

function Render({ seed, ...rest }: StoryArgs) {
  const client = createMockTagsClient({ initial: seed });
  return (
    <TagEditor
      {...rest}
      client={client}
      onChange={(tags) => console.log('tags changed:', tags)}
    />
  );
}

export const Empty: Story<StoryArgs> = (args) => <Render {...args} />;
Empty.args = { subject: SUBJECT, seed: [] };

export const WithExistingTags: Story<StoryArgs> = (args) => <Render {...args} />;
WithExistingTags.args = {
  subject: SUBJECT,
  seed: [
    tag({ uuid: 't1', purpose: 'env', value: 'production', color: '#dc2626' }),
    tag({ uuid: 't2', purpose: 'team', value: 'platform', color: '#2563eb' }),
  ],
};

export const FixedPurpose: Story<StoryArgs> = (args) => <Render {...args} />;
FixedPurpose.args = {
  subject: SUBJECT,
  purposes: ['env'],
  seed: [tag({ uuid: 't1', purpose: 'env', value: 'production', color: '#dc2626' })],
};

export const SelectPurpose: Story<StoryArgs> = (args) => <Render {...args} />;
SelectPurpose.args = {
  subject: SUBJECT,
  purposes: ['env', 'team', 'priority'],
  seed: [
    tag({ uuid: 't1', purpose: 'env', value: 'production', color: '#dc2626' }),
    tag({ uuid: 't2', purpose: 'priority', value: 'p0' }),
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
