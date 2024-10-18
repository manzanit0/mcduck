BEGIN;

CREATE TYPE receipt_status AS ENUM (
  'uploaded', 
  'failed_preprocessing',
  'pending_review', 
  'reviewed'
);

ALTER TABLE receipts
ADD COLUMN status receipt_status NOT NULL DEFAULT 'uploaded';

-- Update existing data.
UPDATE receipts SET status='reviewed' WHERE pending_review=false;
UPDATE receipts SET status='pending_review' WHERE pending_review=true;

COMMIT;
