<!-- Generated with Stardoc: http://skydoc.bazel.build -->

Example rules to show package naming techniques.

<a id="name_part_from_command_line"></a>

## name_part_from_command_line

<pre>
load("@gotopt2//build:my_package_name.bzl", "name_part_from_command_line")

name_part_from_command_line(<a href="#name_part_from_command_line-name">name</a>)
</pre>



**ATTRIBUTES**


| Name  | Description | Type | Mandatory | Default |
| :------------- | :------------- | :------------- | :------------- | :------------- |
| <a id="name_part_from_command_line-name"></a>name |  A unique name for this target.   | <a href="https://bazel.build/concepts/labels#target-names">Name</a> | required |  |


