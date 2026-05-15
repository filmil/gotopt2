def gotopt2_generate_bash(name, src, out):
    """Generates a bash parser script from a YAML configuration using gotopt2-generator.

    Args:
        name: A unique name for this target.
        src: The input YAML file containing the configuration.
        out: The output bash script filename to generate.
    """
    native.genrule(
        name = name,
        srcs = [src],
        outs = [out],
        tools = ["@multitool//tools/gotopt2-generator"],
        cmd = "$(location @multitool//tools/gotopt2-generator) < $(location %s) > $@" % src,
    )
