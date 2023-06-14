#resource "aws_route53_record" "console_record" {
#  zone_id = "Z04527933M5SONRD175ZJ"
#  name    = "aabg.io"
#  type    = "A"
#
#  alias {
#    name                   = "foo"
#    zone_id                = "bar"
#    evaluate_target_health = true
#  }
#
#  alias {
#    name                   = "foo"
#    zone_id                = "bar"
#    evaluate_target_health = true
#  }
#}