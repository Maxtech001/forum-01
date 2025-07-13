provider "aws" {
  region = var.aws_region
}

# IAM Role for SSM
resource "aws_iam_role" "ssm_role" {
  name = "forum-app-ssm-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = {
        Service = "ec2.amazonaws.com"
      },
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ssm_attach" {
  role       = aws_iam_role.ssm_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "ssm_instance_profile" {
  name = "forum-app-ssm-instance-profile"
  role = aws_iam_role.ssm_role.name
}

# Security Group
resource "aws_security_group" "forum_sg" {
  name        = "forum-sg"
  description = "Allow HTTP, HTTPS, SSH and Kubernetes traffic"

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "K3s API"
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Control Plane
resource "aws_instance" "control_plane" {
  ami                         = var.ami_id
  instance_type               = var.instance_type
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.ssm_instance_profile.name
  vpc_security_group_ids      = [aws_security_group.forum_sg.id]

  user_data = <<-EOF
              #!/bin/bash
              set -e
              setenforce 0 || true
              sed -i 's/^SELINUX=.*/SELINUX=permissive/' /etc/selinux/config

              curl -sfL https://get.k3s.io | sh -s - --write-kubeconfig-mode 644

              TOKEN=$(sudo cat /var/lib/rancher/k3s/server/node-token)
              echo $TOKEN > /home/ubuntu/k3s_token.txt
              chown ubuntu:ubuntu /home/ubuntu/k3s_token.txt

              # Store token in SSM
              REGION="${var.aws_region}"
              yum install -y awscli
              aws ssm put-parameter --name "/k3s/token" --value "$TOKEN" --type "SecureString" --overwrite --region $REGION
              EOF

  tags = {
    Name = "Forum-K8s-ControlPlane"
  }
}

# Worker Nodes
resource "aws_instance" "worker_nodes" {
  count                       = 2
  ami                         = var.ami_id
  instance_type               = var.instance_type
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.ssm_instance_profile.name
  vpc_security_group_ids      = [aws_security_group.forum_sg.id]

  user_data = <<-EOF
              #!/bin/bash
              set -e

              # Install K3s
              curl -sfL https://get.k3s.io | sh -s - --write-kubeconfig-mode 644

              # Get K3s token
              TOKEN=$(sudo cat /var/lib/rancher/k3s/server/node-token)
              echo $TOKEN > /home/ubuntu/k3s_token.txt
              chown ubuntu:ubuntu /home/ubuntu/k3s_token.txt

              # Store token in SSM
              apt-get update -y && apt-get install -y awscli
              aws ssm put-parameter --name "/k3s/token" --value "$TOKEN" --type "SecureString" --overwrite --region ${var.aws_region}
              EOF

  tags = {
    Name = "Forum-K8s-Worker-${count.index + 1}"
  }
}

# EIP for Control Plane
resource "aws_eip" "control_plane_eip" {}

resource "aws_eip_association" "control_plane_eip_assoc" {
  instance_id   = aws_instance.control_plane.id
  allocation_id = aws_eip.control_plane_eip.id
}

# EIPs for Worker Nodes
resource "aws_eip" "worker_eips" {
  count = 2
}

resource "aws_eip_association" "worker_eip_assoc" {
  count         = 2
  instance_id   = aws_instance.worker_nodes[count.index].id
  allocation_id = aws_eip.worker_eips[count.index].id
}
