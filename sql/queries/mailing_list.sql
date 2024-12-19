-- name: AddEmailToMailingList :one
INSERT INTO mailing_list (
    id, email
) VALUES (
    ?, ?
) 
RETURNING id, email, last_sent_at;
