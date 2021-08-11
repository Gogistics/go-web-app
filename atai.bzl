def _atai_binary_impl(ctx):
    print("analyzing", ctx.label)

atai_binary = rule(
    implementation = _atai_binary_impl,
)

print("bzl file evaluation")