ALTER TABLE tasks ADD COLUMN rrule TEXT DEFAULT NULL;
CREATE INDEX idx_tasks_rrule ON tasks(rrule) WHERE rrule IS NOT NULL;
ALTER TABLE tasks ADD COLUMN due_date DATE;
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
