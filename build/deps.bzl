# Declaring the gotopt2 dependencies
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:utils.bzl", "maybe")

# Include this into any dependencies that want to compile gotopt2 from source.
# This declaration must be updated every time the dependencies in the workspace
# change.
def gotopt2_dependencies():

    excludes = native.existing_rules().keys()

    maybe(
        go_repository,
        name = "com_github_golang_glog",
        commit = "23def4e6c14b",
        importpath = "github.com/golang/glog",
    )

    maybe(
        go_repository,
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        tag = "v0.2.0",
    )

    maybe(
        go_repository,
        name = "in_gopkg_check_v1",
        commit = "20d25e280405",
        importpath = "gopkg.in/check.v1",
    )

    maybe(
        go_repository,
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )

    if "bazel_bats" not in excludes:
        BAZEL_BATS_VERSION = "0.30.0"
        BAZEL_BATS_SHA256 = "9ae647d2db3aa0bd36af84a0a864dce1c4a1c4f7207b240d3a809862944ecb18"

        maybe(
            http_archive,
            name = "bazel_bats",
            strip_prefix = "bazel-bats-%s" % BAZEL_BATS_VERSION,
            urls = [
                "https://github.com/filmil/bazel-bats/archive/refs/tags/v%s.tar.gz" % BAZEL_BATS_VERSION,
            ],
            sha256 = BAZEL_BATS_SHA256,
        )
        # Since we can not "load" here, we must instead inline expand the deps,
        # and call them here. See below.
        bazel_bats_dependencies()

        maybe(
            http_archive,
            name = "rules_pkg",
            urls = [
                "https://mirror.bazel.build/github.com/bazelbuild/rules_pkg/releases/download/0.9.1/rules_pkg-0.9.1.tar.gz",
                "https://github.com/bazelbuild/rules_pkg/releases/download/0.9.1/rules_pkg-0.9.1.tar.gz",
            ],
            sha256 = "8f9ee2dc10c1ae514ee599a8b42ed99fa262b757058f65ad3c384289ff70c4b8",
        )

# Expanded inline @bazel_bats//build:deps.bzl

_BATS_CORE_BUILD = """
sh_library(
    name = "bats_lib",
    srcs = glob(["libexec/**"]),
)

sh_binary(
    name = "bats",
    srcs = ["bin/bats"],
    visibility = ["//visibility:public"],
    deps = [":bats_lib"],
)

sh_library(
    name = "file_setup_teardown_lib",
    srcs = ["test/file_setup_teardown.bats"],
    visibility = ["//visibility:public"],
    data = glob(["test/fixtures/file_setup_teardown/**"]),
)

sh_library(
    name = "junit_formatter_lib",
    srcs = ["test/junit-formatter.bats"],
    visibility = ["//visibility:public"],
    data = glob(["test/fixtures/junit-formatter/**"]),
)

sh_library(
    name = "parallel_lib",
    srcs = ["test/parallel.bats"],
    visibility = ["//visibility:public"],
    data = glob([
        "test/concurrent-coordination.bash",
        "test/fixtures/parallel/**",
    ]),
)

sh_library(
    name = "run_lib",
    srcs = ["test/run.bats"],
    visibility = ["//visibility:public"],
    data = glob(["test/fixtures/run/**"]),
)

sh_library(
    name = "suite_lib",
    srcs = ["test/suite.bats"],
    visibility = ["//visibility:public"],
    data = glob(["test/fixtures/suite/**"]),
)

sh_library(
    name = "test_helper",
    srcs = ["test/test_helper.bash"],
    visibility = ["//visibility:public"],
)

sh_library(
    name = "trace_lib",
    srcs = ["test/trace.bats"],
    visibility = ["//visibility:public"],
    data = glob(["test/fixtures/trace/**"]),
)

exports_files(glob(["test/*.bats"]))
"""

_BATS_ASSERT_BUILD = """
filegroup(
    name = "load_files",
    srcs = [
        "load.bash",
        "src/assert.bash",
    ],
    visibility = ["//visibility:public"],
)
"""

_BATS_SUPPORT_BUILD = """
filegroup(
    name = "load_files",
    srcs = [
        "load.bash",
        "src/error.bash",
        "src/lang.bash",
        "src/output.bash",
    ],
    visibility = ["//visibility:public"],
)
"""

def bazel_bats_dependencies(
    version = "1.7.0",
    sha256 = "ac70c2a153f108b1ac549c2eaa4154dea4a7c1cc421e3352f0ce6ea49435454e",
    bats_assert_version = None,
    bats_assert_sha256 = None,
    bats_support_version = None,
    bats_support_sha256 = None
):
    if not sha256:
        fail("sha256 for bats-core was not supplied.")

    http_archive(
        name = "bats_core",
        build_file_content = _BATS_CORE_BUILD,
        urls = [
            "https://github.com/bats-core/bats-core/archive/refs/tags/v%s.tar.gz" % version,
        ],
        strip_prefix = "bats-core-%s" % version,
        sha256 = sha256,
    )

    if bats_assert_version:
        if not bats_support_version:
            fail("bats-assert version was set, but was missing set version for dependency bats-support.")
        if not bats_assert_sha256:
            fail("sha256 for bats-assert was not supplied.")
        http_archive(
            name = "bats_assert",
            build_file_content = _BATS_ASSERT_BUILD,
            sha256 = bats_assert_sha256,
            strip_prefix = "bats-assert-%s" % bats_assert_version,
            urls = [
                "https://github.com/bats-core/bats-assert/archive/refs/tags/v%s.tar.gz" % bats_assert_version,
            ],
        )
    if bats_support_version:
        if not bats_support_sha256:
            fail("sha256 for bats-support was not supplied.")
        http_archive(
            name = "bats_support",
            build_file_content = _BATS_SUPPORT_BUILD,
            sha256 = bats_support_sha256,
            strip_prefix = "bats-support-%s" % bats_support_version,
            urls = [
                "https://github.com/bats-core/bats-support/archive/refs/tags/v%s.tar.gz" % bats_support_version,
            ],
        )


