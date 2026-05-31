<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="gotopt2_generate_bash"></a>

## gotopt2_generate_bash

<pre>
load("@gotopt2//build:gotopt2.bzl", "gotopt2_generate_bash")

gotopt2_generate_bash(<a href="#gotopt2_generate_bash-name">name</a>, <a href="#gotopt2_generate_bash-src">src</a>, <a href="#gotopt2_generate_bash-out">out</a>, <a href="#gotopt2_generate_bash-generator">generator</a>)
</pre>

Generates a bash parser script from a YAML configuration using gotopt2-generator.

**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="gotopt2_generate_bash-name"></a>name |  A unique name for this target.   |  none |
| <a id="gotopt2_generate_bash-src"></a>src |  The input YAML file containing the configuration.   |  none |
| <a id="gotopt2_generate_bash-out"></a>out |  The output bash script filename to generate.   |  none |
| <a id="gotopt2_generate_bash-generator"></a>generator |  The generator tool to use.   |  `"@multitool//tools/gotopt2-generator"` |


