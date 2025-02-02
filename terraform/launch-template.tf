resource "aws_launch_template" "gitsetrun_lt" {
  name          = "gitsetrun-launch-template"
  image_id      = "ami-00c257e12d6828491"
  instance_type = "t3.micro"
  iam_instance_profile {
    name = aws_iam_instance_profile.gitsetrun_profile.name
  }
  vpc_security_group_ids = [aws_security_group.gitsetrun_sg.id]


  user_data = base64encode(<<EOF
#!/bin/bash
GITHUB_PAT="${var.github_pat}"
REPO_OWNER="${var.github_repo_owner}"
REPO_NAME="${var.github_repo_name}"

GH_RUNNER_TOKEN=$(curl -X POST -H "Authorization: token $GITHUB_PAT" \
    -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/actions/runners/registration-token | jq -r '.token')

mkdir /actions-runner && cd /actions-runner
curl -o actions-runner-linux-x64-2.321.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.321.0/actions-runner-linux-x64-2.321.0.tar.gz
# Optional: Validate the hash
echo "ba46ba7ce3a4d7236b16fbe44419fb453bc08f866b24f04d549ec89f1722a29e  actions-runner-linux-x64-2.321.0.tar.gz" | shasum -a 256 -c
tar xzf ./actions-runner-linux-x64-2.321.0.tar.gz
export RUNNER_ALLOW_RUNASROOT="1"
./config.sh --url https://github.com/vvksh/GitSetRun --token $GH_RUNNER_TOKEN
./run.sh
EOF
  )
}
