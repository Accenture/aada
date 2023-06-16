resource "aws_route53_record" "console_record" {
  zone_id = "Z04527933M5SONRD175ZJ"
  name    = "aada-${data.aws_region.current.name}.aabg.io"
  type    = "CNAME"
  ttl     = 5
  records = [aws_apigatewayv2_domain_name.httpdomain.domain_name_configuration[0].target_domain_name]
}

resource "aws_route53_record" "websocket_record" {
  zone_id = "Z04527933M5SONRD175ZJ"
  name    = "aada-wss-${data.aws_region.current.name}.aabg.io"
  type    = "CNAME"
  ttl     = 5
  records = [aws_apigatewayv2_domain_name.wsdomain.domain_name_configuration[0].target_domain_name]
}
