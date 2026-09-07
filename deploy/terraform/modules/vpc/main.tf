resource "aws_vpc" "main" {
  cidr_block = var.cidr_block
  tags = {
    Name = "ANGEL-VPC"
  }
}

resource "aws_subnet" "main" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}
