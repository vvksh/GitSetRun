resource "aws_security_group" "gitsetrun_sg" {
  name        = "gitsetrun-sg"
  description = "Security group for gitsetrun"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"] # Allow all outbound traffic
  }
}
