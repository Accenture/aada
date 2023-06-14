resource "aws_acm_certificate" "http_cert" {
  domain_name = "aabg.io"
  validation_method = "DNS"
  subject_alternative_names = ["aabg.io", "wss.aabg.io"]
}
