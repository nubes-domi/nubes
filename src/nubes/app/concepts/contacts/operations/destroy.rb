module Contacts::Operations
  class Destroy < DestroyOperation
    def load(ctx, id:, **)
      ctx["model"] = ctx["current_user"].contacts.find_by(id:)
    end
  end
end
