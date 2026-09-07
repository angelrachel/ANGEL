resource "aws_instance" "lab" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  subnet_id     = var.subnet_id
  tags = {
    Name = "ANGEL-LAB"
  }
}
