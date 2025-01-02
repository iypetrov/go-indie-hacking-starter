resource "aws_secretsmanager_secret" "go_indie_hacking_starter_addr" {
  name            = "go_indie_hacking_starter_addr"
}

resource "aws_secretsmanager_secret_version" "go_indie_hacking_starter_addr_version" {
  secret_id     = aws_secretsmanager_secret.go_indie_hacking_starter_addr.id
  secret_string = var.go_indie_hacking_starter_addr
}

resource "aws_secretsmanager_secret" "go_indie_hacking_starter_db_file" {
  name            = "go_indie_hacking_starter_db_file"
}

resource "aws_secretsmanager_secret_version" "go_indie_hacking_starter_db_file_version" {
  secret_id     = aws_secretsmanager_secret.go_indie_hacking_starter_db_file.id
  secret_string = var.go_indie_hacking_starter_db_file
}