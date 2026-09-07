resource "aws_security_group" "angel" {
  name        = "ANGEL-SEC"
  description = "ANGEL OPSEC"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
