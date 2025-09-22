
def run_gotopt2_sample(name):
    _filename = "{}.txt".format(name)
    native.genrule(
        name = name,
        srcs = [],
        outs = ["{}.txt".format(name)],
        tools = [
            Label("@gotopt2//build:example"),
        ],
        cmd = """$(location @gotopt2//build:example) --help &> $(location {out})"""
            .format(out=_filename),
    )
