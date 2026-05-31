def gotopt2_generate_bash(name, src, out, generator = "@multitool//tools/gotopt2-generator"):
    """Generates a bash parser script from a YAML configuration using gotopt2-generator.

    Args:
        name: A unique name for this target.
        src: The input YAML file containing the configuration.
        out: The output bash script filename to generate.
        generator: The generator tool to use.
    """
    native.genrule(
        name = name,
        srcs = [src],
        outs = [out],
        tools = [generator],
        cmd = "$(location %s) < $(location %s) > $@" % (generator, src),
    )
