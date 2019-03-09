# `gotopt2`: a self-contained shell flags or options parser, written in go

`gotopt2` is a program that outputs its command line arguments as a
snippet of shell script that can be readily evaluated. You can use it to 
parse command line options in your shell script instead of rolling your own,
or using `getopt` or similar.

## Quick example

```console
gotopt2 <<EOF -a -b=foo -c=10 --name value arg1 arg2
- flag:
  short: "a"
  type: bool
  help: "Some help text"
EOF
# gotopt2:generated:begin
readonly gotopt2_a=true
readonly gotopt2_b="foo"
readonly gotopt2_c=10
readonly gotopt2_name="value"
readonly gotopt2_rest=("arg1" "arg2")
# gotopt2:generated:end
```

# Use Case

Though it's 2019, I ofter find myself needing to write `bash` scripts.  As you
may be aware, there isn't really a canned way in `bash` to parse command line
options: you can either roll your own, or you can rely on a preexisting solution
like GNU `getopt`.

As an alternative you could pass options easily in environment variables, but
that ends up being spooky action at a distance when you have your flags passed
through multiple levels of scripts, all alike.

Rolling your own means either you write custom parsing code in each script.
Or, you make it a library, in which case you have to worry about how you
package and load the library in your script when you want to use it.  All
doable, just doesn't feel very efficient when alternatives exist.

If you don't want to roll your own, you could use GNU `getopt` for example.
However, then you need to make sure that you have exactly the version of
`getopt` you need on the target system.  Ensuring this is the job of GNU
`autotools` but as soon as you touch `autotools` it is probably an overkill.
Remember, the only thing you actually wanted is to parse some command line
options. And if you are on OSX, who knows which `getopt` you will be up
against. This was fine in the eighties, but not today. *There must be an easier
way.*

There are libraries like `argbash.io`.  Which I liked very much untl I learned
that (1) I need `m4` to use it and (2) I need a Makefile to generate the actual
running script from my code. At that point it becomes obvious to me that a small
binary works better.

And even if `getopt` fits your bill, you still need to figure out its arcane
flag parsing syntax. Again this was fine in the 1970's, and even desirable as
few computers had actual monitors but printed output on paper instead and pithy
was king.

So I set out to write something to improve on the situation.

# Requirements

*Option parsing should be embeddable in `bash` scripts.*

This is pretty much the main functional requirement; we need to convert command
line arguments into some bash code, and evaluate that.

*Option parsing should be easily portable.*

We get portability by making the program small and self-contained.  So if is
not available as a binary for your platform you can compile it yourself on the
spot.

*Option parsing configuration should be easily readable.*

While I appreciate compact notations as much as the next person, I appreciate
being able to maintain my scripts more. And I appreciate even more the ability
to maintain scripts that *someone else* wrote. For this reason, configuration
should be in an easily understandable, preferably self-documenting forms.

*Option parsing approach should be contemporary.*

This means for example that you get `--help` for free.  And that the help
text is auto-generated from the information you pass at configuration time.
And that both long and short option names are supported.

# Q&A

## Why is it named gotopt2?

I thought I was being original by riffing on the very well known name "getopt"
and weaving in the name of the language the program is written in at the same
time.

However, I was not the first who came to that idea, as you can readily see at
https://github.com/akutz/gotopt.  But the intention of that `gotopt` is to
replicate the functionality of getopt in go (why would you if you had a chance
to redo it?!), so I thought it appropriate to name this `getopt2`.  Looking
forward to `getopt3` if you feel so inclined.

## What's wrong with getopt?

Nothing, if it fulfills your use case.

I have, however, found it very arcane to use, and quite difficult to install
where I needed it.  And there are two versions of it. 

In contrast, one self-contained binary that you can build on the spot is very
easy to carry around and include in your own code, so I opted for that.

## What's wrong with `gotopt`?

Nothing, if it fulfills your use case.

It's reimplementing `getopt`.  I don't see why one would want to do that if
given the opportunity to implement a user-friendlier approach.

## What's wrong with argbash (https://argbash.io)?

You need to carry `m4` around with you.  Which, given that it's 2019 through my
window when I write this, is a horrible predicament.

You also need to have a `Makefile` that builds your final script, which is not
necessarily something you'd want to do.  All else being equal, a small go
program is in my experience a much more robust building block than an untested
and  sprawling `bash` string parsing library.

## Why does it matter which programming language it is written in?

You are right, it doesn't really matter.  However, for some reason in my head
go programs are associated with small static binaries that you can build very
easily and deploy alongside anything you do.  So it's a signal that it is 
something small, simple and portable.

Simplicity and portability become important in a world where packaging is a
solved problem: you can easily build a container image that has this as the
only one additional component.  Compare to say needing to install the whole GNU
`autotools` package if you want to use `argbash`. Or, compare to `getopt` where
the expectation is that your system has one, but it's never the version you
need.
