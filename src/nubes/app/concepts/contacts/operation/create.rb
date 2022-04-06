module Contacts
  module Operation
    class Create < Trailblazer::Operation
      step :create

      def create(ctx, attributes:, current_user:, **)
        ctx["model"] = current_user.contacts.new(attributes)
        ctx["model"].save
      end
    end
  end
end
