# Declaring the gotopt2 dependencies
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_bats//:deps.bzl", "bazel_bats_dependencies")

# Include this into any dependencies that want to compile gotopt2 from source.
# This declaration must be updated every time the dependencies in the workspace
# change.
def gotopt2_dependencies():
    go_repository(
        name = "com_github_golang_glog",
        commit = "23def4e6c14b",
        importpath = "github.com/golang/glog",
    )

    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        tag = "v0.2.0",
    )

    go_repository(
        name = "in_gopkg_check_v1",
        commit = "20d25e280405",
        importpath = "gopkg.in/check.v1",
    )

    go_repository(
        name = "in_gopkg_yaml_v2",
        importpath = "gopkg.in/yaml.v2",
        tag = "v2.2.2",
    )

    git_repository(
        name = "bazel_bats",
        remote = "https://github.com/filmil/bazel-bats",
        commit = "78da0822ea339bd0292b5cc0b5de6930d91b3254",
        shallow_since = "1569564445 -0700",
    )
    bazel_bats_dependencies()

