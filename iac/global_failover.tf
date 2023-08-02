locals {
  us_east_1_domain = regex("https?://([a-z0-9.-]+)/?", module.aada_us_east_1.http_function_url)[0]
  us_west_1_domain = regex("https?://([a-z0-9.-]+)/?", module.aada_us_west_1.http_function_url)[0]
}

resource "aws_cloudfront_distribution" "aada" {
  enabled         = true
  is_ipv6_enabled = true

  origin {
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
    domain_name = local.us_east_1_domain
    origin_id   = "us-east-1"
  }

  origin {
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
    domain_name = local.us_west_1_domain
    origin_id   = "us-west-1"
  }

  origin_group {
    origin_id = "regionalEndpoints"

    failover_criteria {
      status_codes = [500, 502, 503, 504]
    }
    member {
      origin_id = "us-east-1"
    }
    member {
      origin_id = "us-west-1"
    }
  }

  default_cache_behavior {
    target_origin_id       = "us-east-1"
    viewer_protocol_policy = "redirect-to-https"
    cache_policy_id        = "4135ea2d-6df8-44a3-9df3-4b5a84be39ad"
    allowed_methods        = ["HEAD", "DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"]
    cached_methods         = ["GET", "HEAD"]
  }

  viewer_certificate {
    acm_certificate_arn            = "arn:aws:acm:us-east-1:464079168809:certificate/f31a69cb-b236-4249-b4c1-0d4d9bcaae85"
    minimum_protocol_version       = "TLSv1.2_2021"
    ssl_support_method             = "sni-only"
    cloudfront_default_certificate = false
  }

  restrictions {
    geo_restriction {
      restriction_type = "blacklist"
      locations        = ["CU", "IR", "CN", "RU", "KP", "SY"]
    }
  }
}

resource "aws_route53_record" "apex" {
  name    = "aabg.io"
  type    = "A"
  zone_id = "Z04527933M5SONRD175ZJ"

  alias {
    name                   = "d-j2zf4098xg.execute-api.us-east-1.amazonaws.com"
    zone_id                = "Z1UJRXOUMOOFQ8"
    evaluate_target_health = false
  }
}
