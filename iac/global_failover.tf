locals {
  us_east_1_domain      = regex("https?://([a-z0-9.-]+)/?", module.aada_us_east_1.http_function_url)[0]
  us_west_1_domain      = regex("https?://([a-z0-9.-]+)/?", module.aada_us_west_1.http_function_url)[0]
  eu_central_1_domain   = regex("https?://([a-z0-9.-]+)/?", module.aada_eu_central_1.http_function_url)[0]
  ap_southeast_2_domain = regex("https?://([a-z0-9.-]+)/?", module.aada_ap_southeast_2.http_function_url)[0]
}

resource "aws_acm_certificate" "aada" {
  domain_name               = "aabg.io"
  subject_alternative_names = ["wss.aabg.io"]
  key_algorithm             = "EC_prime256v1"
  validation_method         = "DNS"
}

resource "aws_cloudfront_distribution" "aada" {
  enabled         = true
  is_ipv6_enabled = true

  aliases = ["aabg.io"]

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

  origin {
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
    domain_name = local.eu_central_1_domain
    origin_id   = "eu-central-1"
  }

  origin {
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
    domain_name = local.ap_southeast_2_domain
    origin_id   = "ap-southeast-2"
  }

  origin_group {
    origin_id = "usEndpoints"

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

  origin_group {
    origin_id = "alternateEndpoints"

    failover_criteria {
      status_codes = [500, 502, 503, 504]
    }
    member {
      origin_id = "eu-central-1"
    }
    member {
      origin_id = "ap-southeast-2"
    }
  }

  default_cache_behavior {
    target_origin_id       = "usEndpoints"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["HEAD", "GET"]
    cached_methods         = ["GET", "HEAD"]

    forwarded_values {
      query_string = true
      cookies {
        forward = "none"
      }
    }
  }

  viewer_certificate {
    acm_certificate_arn            = aws_acm_certificate.aada.arn
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
    name                   = aws_cloudfront_distribution.aada.domain_name
    zone_id                = aws_cloudfront_distribution.aada.hosted_zone_id
    evaluate_target_health = false
  }
}

moved {
  from = aws_route53_record.wss_east_1
  to   = aws_route53_record.wss_us_east_1
}

resource "aws_route53_record" "wss_us_east_1" {
  name           = "wss.aabg.io"
  type           = "CNAME"
  zone_id        = "Z04527933M5SONRD175ZJ"
  ttl            = 300
  records        = [module.aada_us_east_1.ws_domain_name]
  set_identifier = "wss_us_east_1"

  geolocation_routing_policy {
    continent = "NA"
  }
}

resource "aws_route53_record" "wss_us_west_1" {
  name           = "wss.aabg.io"
  type           = "CNAME"
  zone_id        = "Z04527933M5SONRD175ZJ"
  ttl            = 300
  records        = [module.aada_us_west_1.ws_domain_name]
  set_identifier = "wss_us_west_1"

  geolocation_routing_policy {
    country = "*"
  }
}

resource "aws_route53_record" "wss_eu_central_1" {
  name           = "wss.aabg.io"
  type           = "CNAME"
  zone_id        = "Z04527933M5SONRD175ZJ"
  ttl            = 300
  records        = [module.aada_eu_central_1.ws_domain_name]
  set_identifier = "wss_eu_central_1"

  geolocation_routing_policy {
    continent = "EU"
  }
}

resource "aws_route53_record" "wss_ap_southeast_2" {
  name           = "wss.aabg.io"
  type           = "CNAME"
  zone_id        = "Z04527933M5SONRD175ZJ"
  ttl            = 300
  records        = [module.aada_ap_southeast_2.ws_domain_name]
  set_identifier = "wss_ap_southeast_2"

  geolocation_routing_policy {
    continent = "OC"
  }
}
