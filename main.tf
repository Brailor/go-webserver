terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  profile = "default"
  region  = "eu-central-1"
}

resource "aws_instance" "go_server" {
  ami           = "ami-07df274a488ca9195"
  instance_type = "t2.micro"
  user_data     = file("./app.yml")

  tags = {
    Name = "ExampleAppServerInstance"
  }
}

resource "aws_security_group" "web-sg" {
  name = "my-sg"
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

output "web-address" {
  value = "${aws_instance.go_server.public_dns}:8080"
}