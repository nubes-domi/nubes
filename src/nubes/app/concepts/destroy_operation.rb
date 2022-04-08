class DestroyOperation < Trailblazer::Operation
  step :load
  fail :not_found, fail_fast: true
  step :destroy
  fail :extract_errors

  def not_found(ctx, **)
    ctx["errors"] = "not_found"
  end

  def destroy(ctx, **)
    ctx["model"].destroy
  end

  def extract_errors(ctx, **)
    ctx["errors"] = ctx["model"].errors
  end
end
