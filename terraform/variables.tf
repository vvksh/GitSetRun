variable "github_pat" {
  description = "GitHub Personal Access"
  type        = string
  sensitive   = true
}

variable "github_repo_owner" {
  description = "GitHub Repository Owner"
  type        = string
}

variable "github_repo_name" {
  description = "GitHub Repository Name"
  type        = string
}
