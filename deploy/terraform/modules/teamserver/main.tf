resource "aws_instance" "teamserver" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  tags = {
    Name = "ANGEL-TEAMSERVER"
  }
}
