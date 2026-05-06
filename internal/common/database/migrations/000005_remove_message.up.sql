ALTER TABLE messages ADD COLUMN conversation_id INT;

UPDATE messages m
SET conversation_id = c.conversation_id
FROM conversations c
WHERE c.conversation_type = 2 AND c.reference_id = m.group_id;


ALTER TABLE messages ALTER COLUMN conversation_id SET NOT NULL;
ALTER TABLE messages DROP COLUMN group_id;
ALTER TABLE messages ADD CONSTRAINT fk_message_conversation
    FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id);


DROP TABLE IF EXISTS direct_messages;

DROP INDEX IF EXISTS idx_messages_group_id;
DROP INDEX IF EXISTS idx_messages_group_cursor;
CREATE INDEX IF NOT EXISTS idx_messages_conversation_cursor
    ON messages (conversation_id, message_created_at DESC, message_id DESC)
    INCLUDE (message_content, attachments, message_type, reply_to, message_deleted_at)
    WHERE message_deleted_at IS NULL;
