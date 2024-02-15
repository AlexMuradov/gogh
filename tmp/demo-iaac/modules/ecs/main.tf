resource "aws_security_group" "project" {
  name        = "${var.projectname}-lbs-sg"
  description = "${var.projectname} lbs sg"
  vpc_id      = var.vpc_id

  ingress {
      description      = "tcp allow"
      from_port        = 0
      to_port          = 0
      protocol         = "tcp"
      cidr_blocks      = ["0.0.0.0/0"]
    }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_security_group" "project-ecs" {
  name        = "${var.projectname}-ecs-sg"
  description = "${var.projectname} ecs sg"
  vpc_id      = var.vpc_id

  ingress {
      description      = "alb allow"
      from_port        = 0
      to_port          = 0
      protocol         = "tcp"
      security_groups  = [aws_security_group.project.id]
    }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

}

resource "aws_security_group_rule" "project" {
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.project.id
}

resource "aws_ecs_cluster" "project" {
  name = "${var.projectname}-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_ecs_cluster_capacity_providers" "project" {
  cluster_name = aws_ecs_cluster.project.name

  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}

resource "aws_cloudwatch_log_group" "log_group" {
  name = var.projectname
}

resource "aws_ecs_task_definition" "project" {
  family = "${var.projectname}-service"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  task_role_arn = var.execrole
  execution_role_arn = var.execrole

  runtime_platform {
    operating_system_family = "LINUX"
  }

  container_definitions = jsonencode([
    {
      name      = "${var.projectname}"
      image     = "${var.rep_url}:${var.image}"
      cpu       = 256
      memory    = 512
    
      execution_role_arn = var.execrole
      essential = true
      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
        }]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.log_group.name
          awslogs-region        = "eu-central-1"
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}

resource "aws_lb" "project" {
  name               = "${var.projectname}-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.project.id]
  subnets            = [for subnet in var.subnet_id : subnet.id]

  enable_deletion_protection = false

}

resource "aws_lb_target_group" "project" {
  name        = "${var.projectname}-tg"
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = var.vpc_id
}

resource "aws_lb_listener" "project" {
  load_balancer_arn = aws_lb.project.arn
  port              = "80"
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.project.arn
  }
}

resource "aws_ecs_service" "project" {
  name            = "${var.projectname}-svc"
  cluster         = aws_ecs_cluster.project.id
  task_definition = aws_ecs_task_definition.project.arn
  desired_count   = 1
  depends_on = [aws_security_group.project-ecs]
  lifecycle {
    ignore_changes = [
      capacity_provider_strategy
    ]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.project.id
    container_name   = "${var.projectname}"
    container_port   = 80
  }

  network_configuration {
    subnets = [for subnet in var.subnet_id : subnet.id]
    security_groups = [aws_security_group.project-ecs.id]
    assign_public_ip = true
  }

}