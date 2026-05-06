DO $$
DECLARE
    constr_name text;
BEGIN
    SELECT conname INTO constr_name
    FROM pg_constraint
    WHERE conrelid = 'users'::regclass
      AND contype = 'c'
      AND conkey @> (SELECT array_agg(attnum) FROM pg_attribute WHERE attrelid = 'users'::regclass AND attname = 'user_status');

    IF constr_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE users DROP CONSTRAINT %I', constr_name);
    END IF;
END $$;

ALTER TABLE users ADD CONSTRAINT chk_user_status CHECK (user_status IN (1,2,3));


ALTER TABLE conversations DROP CONSTRAINT IF EXISTS chk_conversation_type;
ALTER TABLE conversations ADD CONSTRAINT chk_conversation_type CHECK (conversation_type IN (1,2));

ALTER TABLE conversations DROP CONSTRAINT IF EXISTS uq_conversation_type_ref;
ALTER TABLE conversations ADD CONSTRAINT uq_conversation_type_ref UNIQUE (conversation_type, reference_id);


DROP INDEX IF EXISTS uq_active_group_member;
CREATE UNIQUE INDEX IF NOT EXISTS uq_active_group_member ON group_members (group_id, user_id) WHERE left_at IS NULL;


ALTER TABLE conversation_participants DROP CONSTRAINT IF EXISTS conversation_participants_conversation_id_user_id_key;
DROP INDEX IF EXISTS uq_active_conv_participant;
CREATE UNIQUE INDEX IF NOT EXISTS uq_active_conv_participant ON conversation_participants (conversation_id, user_id) WHERE left_at IS NULL;


CREATE INDEX IF NOT EXISTS idx_messages_group_cursor
    ON messages (group_id, message_created_at DESC, message_id DESC)
    INCLUDE (message_content, attachments, message_type, reply_to, message_deleted_at)
    WHERE message_deleted_at IS NULL;

UPDATE messages
SET reply_to = NULL
WHERE reply_to IS NOT NULL
    AND reply_to NOT IN (SELECT message_id FROM messages);
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_reply_to;
ALTER TABLE messages ADD CONSTRAINT fk_reply_to FOREIGN KEY (reply_to) REFERENCES messages(message_id);


ALTER TABLE group_members DROP CONSTRAINT fk_group;
ALTER TABLE group_members ADD CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE;

ALTER TABLE messages DROP CONSTRAINT fk_group;
ALTER TABLE messages ADD CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE;


ALTER TABLE conversation_participants DROP CONSTRAINT fk_conversation;
ALTER TABLE conversation_participants ADD CONSTRAINT fk_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id) ON DELETE CASCADE;

CREATE SEQUENCE IF NOT EXISTS dm_reference_seq START 1;
