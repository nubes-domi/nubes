module Contacts::Operations
  class Update < Trailblazer::Operation
    step :load
    fail :not_found, fail_fast: true
    step :update
    fail :extract_errors

    def load(ctx, id:, **)
      ctx["model"] = ctx["current_user"].contacts.find_by(id:)
    end

    def not_found(ctx, **)
      ctx["errors"] = "not_found"
    end

    def update(ctx, params:, **)
      ctx["model"].update(params)
    end

    def extract_errors(ctx, **)
      ctx["errors"] = ctx["model"].errors
    end
  end
end
