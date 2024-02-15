resource "aws_codecommit_repository" "project" {
  repository_name = "project"
  description     = "project"
}

resource "aws_iam_role" "project" {
  name                = "DevToolsRole"
  managed_policy_arns = [aws_iam_policy.policy_one.arn, aws_iam_policy.policy_two.arn]
}

resource "aws_s3_bucket" "project" {
  bucket = "projectapp-log-bucket"
}

resource "aws_s3_bucket_acl" "project" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}

resource "aws_s3_bucket_versioning" "versioning_project" {
  bucket = aws_s3_bucket.example.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_codebuild_project" "project" {
  name          = "project"
  description   = "project codebuild project"
  build_timeout = "5"
  service_role  = aws_iam_role.example.arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  cache {
    type     = "S3"
    location = aws_s3_bucket.example.bucket
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/standard:1.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"

    environment_variable {
      name  = "SOME_KEY1"
      value = "SOME_VALUE1"
    }

    environment_variable {
      name  = "SOME_KEY2"
      value = "SOME_VALUE2"
      type  = "PARAMETER_STORE"
    }
  }

  logs_config {
    cloudwatch_logs {
      group_name  = "log-group"
      stream_name = "log-stream"
    }

    s3_logs {
      status   = "ENABLED"
      location = "${aws_s3_bucket.example.id}/build-log"
    }
  }

  source {
    type            = "CODECOMMIT"
    location        = "https://git-codecommit.eu-central-1.amazonaws.com/v1/repos/${aws_codecommit_repository.project.repository_name}"
    git_clone_depth = 1

    git_submodules_config {
      fetch_submodules = true
    }
  }

  source_version = "master"

  vpc_config {
    vpc_id = var.vpc_id

    subnets = [
      var.subnet_id
    ]

    security_group_ids = [
      var.security_group_id
    ]
  }

  tags = {
    Environment = "Test"
  }
}