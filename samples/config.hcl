root = "./samples"
out_directory = "./samples/out"

directory "user_action" {
  include = [
    ".*"]
  exclude = []
  column "a" {
    type = "Int"
    path = "a.s.a"
  }
  column "b" {
    type = "String"
    path = "a.s.b"
    skip = true
  }
  column "c" {
    type = "String"
    path = "a.s.c"
  }
  column "d" {
    type = "Boolean"
    path = "a.s.d"
  }
  column "d" {
    type = "String"
    path = "a.s.ip"
    default = "127.0.0.1"
  }

  column "include" {
    type = "Int"
    path = "an_int"
    includes = [
      1,
      2,
      3]
  }
  column "excludes" {
    type = "String"
    path = "a_string"
    excludes = [
      "a"]
  }

  additional_column "type" {
    type = "String"
    path = "a.type"
    default = "PING"
  }
}


directory "user_info" {
  include = [
    ".*"]
  exclude = []
  separator = "\t"
  column "a" {
    type = "Int"
    path = "a.s.a"
  }
  column "b" {
    type = "String"
    path = "a.s.b"
    skip = true
  }
  column "c" {
    type = "String"
    path = "a.s.c"
  }
  column "d" {
    type = "Boolean"
    path = "a.s.d"
  }
  column "d" {
    type = "String"
    path = "a.s.ip"
    default = "127.0.0.1"
  }
  additional_column "type" {
    type = "String"
    path = "a.type"
    default = "PING"
  }
}
