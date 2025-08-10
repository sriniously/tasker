CREATE TABLE todo_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    todo_id UUID NOT NULL REFERENCES todos ON DELETE CASCADE,
    name TEXT NOT NULL,
    uploaded_by TEXT NOT NULL,
    download_key TEXT NOT NULL,
    file_size BIGINT,
    mime_type TEXT
);

-- Indexes for todo_attachments
CREATE INDEX idx_todo_attachments_todo_id ON todo_attachments(todo_id);
CREATE INDEX idx_todo_attachments_uploaded_by ON todo_attachments(uploaded_by);

-- Updated at trigger for todo_attachments
CREATE TRIGGER set_updated_at_todo_attachments
    BEFORE UPDATE ON todo_attachments
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();