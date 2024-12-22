load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")


### BEGIN: Go with Gazelle
load("//build:deps0.bzl", "gotopt2_deps0")
gotopt2_deps0()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
go_rules_dependencies()
go_register_toolchains(version = "1.20.7")
gazelle_dependencies()
### END: Go with Gazelle


# Ostensibly this is the only thing needed for gotopt2 to work.
load("//build:deps.bzl", "gotopt2_dependencies")
gotopt2_dependencies()

load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")
rules_pkg_dependencies()
