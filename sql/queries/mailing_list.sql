-- name: AddEmailToMailingList :one
INSERT INTO mailing_list (
    id, email, created_at
) VALUES (
    ?, ?, ?
) 
RETURNING id, email, last_sent_at, created_at;
