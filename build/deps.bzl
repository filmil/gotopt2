# Declaring the gotopt2 dependencies
load("@bazel_bats//:deps.bzl", "bazel_bats_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Include this into any dependencies that want to compile gotopt2 from source.
# This declaration must be updated every time the dependencies in the workspace
# change.
def gotopt2_dependencies():

    excludes = native.existing_rules().keys()

    if "com_github_golang_glog" not in excludes:
        go_repository(
            name = "com_github_golang_glog",
            commit = "23def4e6c14b",
            importpath = "github.com/golang/glog",
        )

    if "com_github_google_go_cmp" not in excludes:
        go_repository(
            name = "com_github_google_go_cmp",
            importpath = "github.com/google/go-cmp",
            tag = "v0.2.0",
        )

    if "in_gopkg_check_v1" not in excludes:
        go_repository(
            name = "in_gopkg_check_v1",
            commit = "20d25e280405",
            importpath = "gopkg.in/check.v1",
        )

    if "in_gopkg_yaml_v2" not in excludes:
        go_repository(
            name = "in_gopkg_yaml_v2",
            importpath = "gopkg.in/yaml.v2",
            tag = "v2.2.2",
        )

    if "bazel_bats" not in excludes:
        BAZEL_BATS_VERSION = "0.30.0"
        BAZEL_BATS_SHA256 = "9ae647d2db3aa0bd36af84a0a864dce1c4a1c4f7207b240d3a809862944ecb18"

        http_archive(
            name = "bazel_bats",
            strip_prefix = "bazel-bats-%s" % BAZEL_BATS_VERSION,
            urls = [
                "https://github.com/filmil/bazel-bats/archive/refs/tags/v%s.tar.gz" % BAZEL_BATS_VERSION,
            ],
            sha256 = BAZEL_BATS_SHA256,
        )

        bazel_bats_dependencies()
