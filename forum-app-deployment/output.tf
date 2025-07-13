output "control_plane_public_ip" {
  value       = aws_instance.control_plane.public_ip
  description = "Public IP of the K8s Control Plane Node"
}

output "worker_node_public_ips" {
  value       = [for instance in aws_instance.worker_nodes : instance.public_ip]
  description = "Public IPs of the K8s Worker Nodes"
}

output "control_plane_static_ip" {
  value       = aws_eip.control_plane_eip.public_ip
  description = "Static Elastic IP for the K8s Control Plane"
}

output "worker_node_static_ips" {
  value       = [for eip in aws_eip.worker_eips : eip.public_ip]
  description = "Static Elastic IPs for the K8s Worker Nodes"
}
