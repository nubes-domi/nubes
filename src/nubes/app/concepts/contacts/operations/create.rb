module Contacts::Operations
  class Create < Trailblazer::Operation
    step :build
    step :save
    fail :extract_errors

    def build(ctx, params:, **)
      ctx["model"] = ctx["current_user"].contacts.build(params)
    end

    def save(ctx, **)
      ctx["model"].save
    end

    def extract_errors(ctx, **)
      ctx["errors"] = ctx["model"].errors
    end
  end
end
